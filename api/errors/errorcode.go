package errors

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
	REQUIRED_PARAMETER_MISSING = ErrorCode{
		Code:    20000,
		Message: "Required parameter missing",
	}
	APP_KEY_NOT_FOUND = ErrorCode{
		Code:    30000,
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
