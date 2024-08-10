// Package middlewares 中间件层
package middlewares

import (
	"VideoWeb/Utilities"
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthAdminCheck 这个中间件验证用户是否为管理员
func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is admin
		auth := c.GetHeader("Authorization")
		userClaim, err := logic.ParseToken(auth)
		if err != nil {
			Utilities.SendErrMsg(c, "middlewares::AuthAdminCheck", http.StatusUnauthorized, "Unauthorized")
			c.Abort() // 中间件验证失败，取消执行后面的流程（关键）
			return
		}
		if userClaim.IsAdmin != 1 {
			Utilities.SendErrMsg(c, "middlewares::AuthAdminCheck", http.StatusUnauthorized, "Not Admin")
			c.Abort() // 中间件验证失败，取消执行后面的流程（关键）
			return
		}
		c.Next()
	}
}
