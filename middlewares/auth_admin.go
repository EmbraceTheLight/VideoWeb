// Package middlewares 中间件层
package middlewares

import (
	"VideoWeb/Utilities"
	"VideoWeb/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthAdminCheck 这个中间件验证用户是否为管理员
func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check if user is admin
		auth := c.GetHeader("Authorization")
		fmt.Println("auth--------->", auth)
		userClaim, err := logic.ParseToken(auth)
		if err != nil {
			Utilities.SendJsonMsg(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort() // TODO:中间件验证失败，取消执行后面的流程（关键）
			return
		}
		if userClaim.IsAdmin != 1 {
			Utilities.SendJsonMsg(c, http.StatusUnauthorized, "Not Admin")
			c.Abort() // TODO:中间件验证失败，取消执行后面的流程（关键）
			return
		}
		c.Next()
	}
}
