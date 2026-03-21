package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(400, gin.H{
		"code": 400,
		"msg":  msg,
		"data": nil,
	})
}

func InternalError(c *gin.Context, msg string) {
	c.JSON(500, gin.H{
		"code": 500,
		"msg":  msg,
		"data": nil,
	})
}

// func Error(c *gin.Context, code int, msg string) {
// 	c.JSON(code, gin.H{
// 		"code": code,
// 		"msg":  msg,
// 		"data": nil,
// 	})
// }
