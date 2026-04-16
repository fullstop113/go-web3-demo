package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
)

type clientWindow struct {
	count     int
	windowEnd time.Time
}

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	clients := make(map[string]clientWindow)

	return func(c *gin.Context) {
		if maxRequests <= 0 || window <= 0 {
			c.Next()
			return
		}

		now := time.Now()
		key := c.ClientIP()

		mu.Lock()
		state, exists := clients[key]
		if !exists || now.After(state.windowEnd) {
			state = clientWindow{count: 0, windowEnd: now.Add(window)}
		}

		state.count++
		clients[key] = state
		allowed := state.count <= maxRequests
		mu.Unlock()

		if !allowed {
			utils.Fail(c, http.StatusTooManyRequests, utils.CodeRateLimit, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
