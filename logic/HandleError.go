package logic

import (
	"VideoWeb/Utilities"
	"github.com/gin-gonic/gin"
)

// AddFuncName 给上下文添加函数名,便于处理错误时记录日志
func AddFuncName(c *gin.Context, funcName string) {
	fname, exist := c.Get("funcName")
	if !exist {
		c.Set("funcName", funcName)
	} else {
		c.Set("funcName", fname.(string)+"::"+funcName)
	}
}

// handleInternalServerError 处理内部错误辅助函数
func handleInternalServerError(c *gin.Context, err error) {
	funcName, _ := c.Get("funcName")
	Utilities.SendErrMsg(c, funcName.(string), 5000, err.Error())
}
