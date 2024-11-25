package utils

import (
	"encoding/json"
	"net/http"
	"online-offline-indicator/types"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, types.ErrResponse{
		CommonResponse: types.CommonResponse{
			Status: types.StatusFail,
			Code:   status,
		},
		Message: err.Error(),
	})
}
