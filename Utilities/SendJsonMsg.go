package Utilities

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendSuccessMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func SendErrMsg(c *gin.Context, funcName string, code int, msg string) {
	SendSuccessMsg(c, code, msg)
	WriteErrLog(funcName, msg)
}
