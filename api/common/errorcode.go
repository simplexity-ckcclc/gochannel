package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

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
	DUPLICATE_APP_KEY = ErrorCode{
		Code:    20101,
		Message: "Duplicate App key",
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

type response struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseJSON(c *gin.Context, httpStatus int, errorCode ErrorCode) {
	c.JSON(httpStatus, errorCode)
}

func ResponseJSONWithExtraMsg(c *gin.Context, httpStatus int, errorCode ErrorCode, extraMsg string) {
	res := &response{
		Code:    errorCode.Code,
		Message: errorCode.Message + " : " + extraMsg,
	}
	c.JSON(httpStatus, res)
}

func ResponseJSONWithData(c *gin.Context, httpStatus int, errorCode ErrorCode, data interface{}) {
	res := &response{
		Code:    errorCode.Code,
		Message: errorCode.Message,
		Data:    data,
	}
	c.JSON(httpStatus, res)
}

func (error ErrorCode) String() string {
	return fmt.Sprintf("Code: %d, Message: %s", error.Code, error.Message)
}
