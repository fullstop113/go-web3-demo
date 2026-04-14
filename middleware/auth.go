package middleware

import (
	"net/http"
	"strings"

	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "user_id"
	ContextUsername = "username"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "缺少Authorization header",
			})
			c.Abort()
			return
		}

		// 格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization header格式错误，应为: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token无效或已过期",
			})
			c.Abort()
			return
		}

		// 注入上下文
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)

		c.Next()
	}
}
