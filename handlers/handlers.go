package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"online-offline-indicator/types"
	"online-offline-indicator/utils"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Handler struct {
	rc *redis.Client
	wg *sync.WaitGroup
}

func NewHandler(rc *redis.Client) *Handler {
	return &Handler{
		rc: rc,
		wg: &sync.WaitGroup{},
	}
}

/******* Handlers *******/

// GET /users?ids=<list of user ids>
func (h *Handler) GetUserStatusByBatch(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if !queryParams.Has("ids") {
		utils.WriteError(w, http.StatusBadRequest, errors.New("query param 'ids' is missing"))
	}
	// Split query param values to fetch list of user IDs
	userIdstr := queryParams.Get("ids")
	userIdList := strings.Split(userIdstr, ",")
	userCh := make(chan string, len(userIdList))
	response := make(map[string]string)
	for _, userId := range userIdList {
		h.wg.Add(1)
		go getUserStatus(h, userCh, userId)
	}
	go func() {
		h.wg.Wait()
		close(userCh)
	}()
	for valueStr := range userCh {
		kvPair := strings.Split(valueStr, "-")
		response[kvPair[0]] = kvPair[1]
	}
	utils.WriteJSON(w, http.StatusOK, types.Response{
		CommonResponse: types.CommonResponse{
			Status: types.StatusPass,
			Code:   http.StatusOK,
		},
		Data: response,
	})
}

// POST /heartbeat
func (h *Handler) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	var payload types.HeartbitRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("payload is invalid"))
	}

	h.rc.Set(r.Context(), payload.Id, "online", 60*time.Second)
	utils.WriteJSON(w, http.StatusOK, types.Response{
		CommonResponse: types.CommonResponse{
			Status: types.StatusPass,
			Code:   http.StatusOK,
		},
		Data: struct{}{},
	})
}

// GET /users/{id}
func (h *Handler) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("user ID is not valid"))
	}

	status, err := h.rc.Get(r.Context(), userId).Result()
	if err != nil {
		status = "offline"
	}

	utils.WriteJSON(w, http.StatusOK, types.Response{
		CommonResponse: types.CommonResponse{
			Status: types.StatusPass,
			Code:   http.StatusOK,
		},
		Data: map[string]string{userId: status},
	})

}

func getUserStatus(h *Handler, ch chan string, id string) {
	defer h.wg.Done()
	if status, err := h.rc.Get(context.Background(), id).Result(); err != nil {
		ch <- strings.Join([]string{id, "offline"}, "-")
	} else {
		ch <- strings.Join([]string{id, status}, "-")
	}
}
