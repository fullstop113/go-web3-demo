package utils

import (
	"github.com/gin-gonic/gin"
)

const (
	CodeOK = 0
	CodeInvalid = 10001
	CodeAuth = 10002
	CodeServer = 10003
)

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: CodeOK,
		Msg: "success",
		Data: data,
	})
}

func Fail(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg: msg,
	})
}