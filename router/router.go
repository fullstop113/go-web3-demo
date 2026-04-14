package router

import (
	"github.com/fullstop113/go-web3-demo/controller"
	"github.com/fullstop113/go-web3-demo/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		userCtrl := controller.NewUserController()

		// 公开路由
		v1.POST("/register", userCtrl.Register)
		v1.POST("/login", userCtrl.Login)

		// 需要认证的路由
		private := v1.Group("")
		private.Use(middleware.JWTAuth())
		{
			private.GET("/user/me", userCtrl.GetCurrentUser)
		}
	}

	return r
}

const Name = "router"
