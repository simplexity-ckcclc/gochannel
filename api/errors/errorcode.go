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
		Code:    20100,
		Message: "App key not found",
	}
	SIG_INVALID = ErrorCode{
		Code:    20200,
		Message: "Signature invalid",
	}
	DUPLICATE_NONCE = ErrorCode{
		Code:    20300,
		Message: "Duplicate nonce",
	}
)

type ErrorCode struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (error ErrorCode) String() string {
	return fmt.Sprintf("Code: %d, Message: %s", error.Code, error.Message)
}
