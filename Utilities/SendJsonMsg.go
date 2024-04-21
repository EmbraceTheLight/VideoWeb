package Utilities

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendErrMsg(c *gin.Context, funcName string, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
	WriteErrLog(funcName, msg)
}
func SendSuccessMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}
