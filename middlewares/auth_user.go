// Package middlewares 中间件层
package middlewares

import (
	"VideoWeb/logic"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthUserCheck 这个中间件验证用户是否为用户
func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check if user is admin
		auth := c.GetHeader("Authorization")
		fmt.Println("auth--------->", auth)
		userClaim, err := logic.ParseToken(auth)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized",
			})
			c.Abort() // TODO:中间件验证失败，取消执行后面的流程（关键）
			return
		}
		if userClaim == nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Not User",
			})
			c.Abort() // TODO:中间件验证失败，取消执行后面的流程（关键）
			return
		}
		c.Set("user", userClaim)
		c.Next()
	}
}
