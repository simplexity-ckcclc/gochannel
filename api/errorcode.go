package api

import "fmt"

var (
	SUCCESS = ErrorCode{
		Code:    0,
		Message: "Success",
	}
	INTERNAL_SERVER_ERROR = ErrorCode{
		Code:    10000,
		Message: "Internal server error",
	}
	APP_KEY_NOT_FOUND = ErrorCode{
		Code:    20000,
		Message: "App key not found",
	}
)

type ErrorCode struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (error ErrorCode) String() string {
	return fmt.Sprintf("Code: %d, Message: %s", error.Code, error.Message)
}
