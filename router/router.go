package router

import (
	"github.com/gin-gonic/gin"
	"github.com/fullstop113/go-web3-demo/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.JWTAuth())
		
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	return r
}

const Name = "router"
