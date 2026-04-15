package middleware

import (
	"net/http"
	"strings"
	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) !=2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}