// Package middlewares 中间件层
package middlewares

import (
	"VideoWeb/logic"
	"github.com/gin-gonic/gin"
)

// CheckIfUserLogin 这个中间件验证用户是否登录成功
func CheckIfUserLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is common user
		auth := c.GetHeader("Authorization")
		//fmt.Println("auth--------->", auth)
		userClaim, err := logic.ParseToken(auth)
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{
			//	"code": http.StatusUnauthorized,
			//	"msg":  "Unauthorized",
			//})
			//c.Abort() // 中间件验证失败，取消执行后面的流程（关键）
			//return
			c.Set("user", nil) //置空user，表明这是未登录用户
			c.Next()
		}
		c.Set("user", userClaim)
		c.Next()
	}
}
