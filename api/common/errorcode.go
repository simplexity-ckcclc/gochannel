package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type E ErrorCode

var (
	Success = ErrorCode{
		Code:    0,
		Message: "Success",
	}
	InternalServerError = ErrorCode{
		Code:    10000,
		Message: "Internal server error",
	}
	RequiredParameterMissing = ErrorCode{
		Code:    20000,
		Message: "Required parameter missing",
	}
	ChannelNotFound = ErrorCode{
		Code:    20100,
		Message: "Channel not found",
	}
	DuplicateChannel = ErrorCode{
		Code:    20101,
		Message: "Duplicate channel",
	}
	AppChannelUnmatched = ErrorCode{
		Code:    20102,
		Message: "App-Channel unmatched",
	}
	ChannelTypeInvalid = ErrorCode{
		Code:    20103,
		Message: "Channel type invalid",
	}
	SigInvalid = ErrorCode{
		Code:    20200,
		Message: "Signature invalid",
	}
	DuplicateNonce = ErrorCode{
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
