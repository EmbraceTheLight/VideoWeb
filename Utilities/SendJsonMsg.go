package Utilities

import (
	"VideoWeb/Utilities/logf"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendJsonMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func SendErrMsg(c *gin.Context, funcName string, code int, msg string) {
	SendJsonMsg(c, code, msg)
	logf.WriteErrLog(funcName, msg)
}
