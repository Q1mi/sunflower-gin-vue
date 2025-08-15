package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 10000, // 业务错误码
	"message": xx,     // 提示信息
	"data": {},    // 数据
}

*/

type ResponseData[T any] struct {
	Code    ResCode `json:"code"`
	Message string  `json:"message"`
	Data    T       `json:"data"`
}

// ResponseError 返回错误信息
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData[any]{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	})
}

// ResponseErrorWithMsg 返回自定义错误信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg string) {
	c.JSON(http.StatusOK, &ResponseData[any]{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// ResponseWithHTTPStatus 返回HTTP状态码和错误
func ResponseErrorWithHTTPStatus(c *gin.Context, status int) {
	c.JSON(status, nil)
}

// ResponseSuccess 返回成功信息
func ResponseSuccess[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, &ResponseData[T]{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	})
}
