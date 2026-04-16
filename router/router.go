package router

import (
	"time"

	"github.com/fullstop113/go-web3-demo/controller"
	"github.com/fullstop113/go-web3-demo/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.RequestLogger(), middleware.Recovery())

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

		article := apiV1.Group("/articles")
		{
			article.GET("", controller.ListArticles)
			article.GET("/:id", controller.GetArticle)
		}

		privateArticle := apiV1.Group("/articles")
		privateArticle.Use(middleware.JWTAuth(), middleware.RateLimit(30, time.Minute))
		{
			privateArticle.POST("", controller.CreateArticle)
			privateArticle.PUT("/:id", controller.UpdateArticle)
			privateArticle.DELETE("/:id", controller.DeleteArticle)
		}
	}
	return r
}

const Name = "router"
