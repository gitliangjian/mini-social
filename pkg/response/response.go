package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	CodeSuccess      = 0
	CodeBadRequest   = 1001
	CodeUnauthorized = 1002
	CodeForbidden    = 1003
	CodeNotFound     = 1004
	CodeInternal     = 1005
)

func JSON(c *gin.Context, httpStatus int, code int, msg string, data any) {
	c.JSON(httpStatus, Resp{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Success(c *gin.Context, data any) {
	JSON(c, http.StatusOK, CodeSuccess, "success", data)
}

func BadRequest(c *gin.Context, msg string) {
	JSON(c, http.StatusBadRequest, CodeBadRequest, msg, nil)
}

func Unauthorized(c *gin.Context, msg string) {
	JSON(c, http.StatusUnauthorized, CodeUnauthorized, msg, nil)
}

func Forbidden(c *gin.Context, msg string) {
	JSON(c, http.StatusForbidden, CodeForbidden, msg, nil)
}

func NotFound(c *gin.Context, msg string) {
	JSON(c, http.StatusNotFound, CodeNotFound, msg, nil)
}

func InternalError(c *gin.Context, msg string) {
	JSON(c, http.StatusInternalServerError, CodeInternal, msg, nil)
}

// func Error(c *gin.Context, code int, msg string) {
// 	c.JSON(code, gin.H{
// 		"code": code,
// 		"msg":  msg,
// 		"data": nil,
// 	})
// }
