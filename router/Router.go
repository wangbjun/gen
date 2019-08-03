package router

import (
	. "gen/controller"
	"gen/middleware"
	"github.com/gin-gonic/gin"
)

func Route(Router *gin.Engine) {
	Router.GET("/", BaseController.Index)

	article := Router.Group("/api/v1/articles").Use(middleware.Auth())
	{
		article.GET("/", ArticleController.ListArticle)
		article.GET("/:id", ArticleController.GetArticle)
		article.POST("/", ArticleController.AddArticle)
		article.DELETE("/:id", ArticleController.DelArticle)
	}

	user := Router.Group("/api/v1/user")
	{
		user.POST("/register", UserController.Register)
		user.POST("/login", UserController.Login)
	}
}
