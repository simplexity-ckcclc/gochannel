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
)

type ErrorCode struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (error ErrorCode) String() string {
	return fmt.Sprintf("Code: %d, Message: %s", error.Code, error.Message)
}
