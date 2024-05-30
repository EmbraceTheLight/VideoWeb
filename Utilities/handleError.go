package Utilities

import (
	"github.com/gin-gonic/gin"
)

// DeferFunc 用于defer函数的调用来处理错误
func DeferFunc(c *gin.Context, err error, funcName string) {
	if err != nil {
		AddFuncName(c, funcName)
	}
}

// AddFuncName 给上下文添加函数名,便于处理错误时记录日志
func AddFuncName(c *gin.Context, funcName string) {
	fname, exist := c.Get("funcName")
	if !exist {
		c.Set("funcName", funcName)
	} else {
		c.Set("funcName", fname.(string)+"::"+funcName)
	}
}

// HandleInternalServerError 处理内部错误辅助函数
func HandleInternalServerError(c *gin.Context, err error) {
	funcName, _ := c.Get("funcName")
	SendErrMsg(c, funcName.(string), 5000, err.Error())
}
