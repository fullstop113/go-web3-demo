package middleware

import (
	"fmt"
	"net/http"

	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				requestID, _ := c.Get("request_id")
				msg := "internal server error"
				if requestID != nil {
					msg = fmt.Sprintf("%s (request_id=%v)", msg, requestID)
				}
				utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, msg)
				c.Abort()
			}
		}()
		c.Next()
	}
}
