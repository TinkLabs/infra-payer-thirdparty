package valoot

import "fmt"

type ErrorResp struct {
	Status     string `json:"status"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e ErrorResp) Error() string {
	return fmt.Sprintf("code: %d, info: %s, status_code: %s", e.Code, e.Message, e.StatusCode)
}
