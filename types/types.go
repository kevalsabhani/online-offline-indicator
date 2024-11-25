package types

const (
	StatusPass = "pass"
	StatusFail = "fail"
)

type HeartbitRequest struct {
	Id string `json:"id"`
}

type CommonResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type ErrResponse struct {
	CommonResponse
	Message string `json:"message"`
}

type Response struct {
	CommonResponse
	Data any `json:"data"`
}
