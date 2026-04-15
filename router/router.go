package router

import (
	"github.com/fullstop113/go-web3-demo/controller"
	"github.com/fullstop113/go-web3-demo/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	apiV1 := r.Group("/api/v1")
	{
		// public routes
		apiV1.POST("/register", controller.Register)
		apiV1.POST("/login", controller.Login)

		// private routes
		private := apiV1.Group("/user")
		private.Use(middleware.JWTAuth())
		{
			private.GET("/getUserInfo", controller.GetUserInfo)
		}
	}
	return r
}

const Name = "router"
