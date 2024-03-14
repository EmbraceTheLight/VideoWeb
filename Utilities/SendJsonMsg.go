package Utilities

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendJsonMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}
