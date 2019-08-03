package router

import (
	. "gen/controller"
	"gen/middleware"
	"github.com/gin-gonic/gin"
)

func Route(Router *gin.Engine) {
	Router.GET("/", BaseController.Index)
	api := Router.Group("/api/v1").Use(middleware.Auth())
	{
		api.GET("/articles", ArticleController.ListArticle)
		api.GET("/articles/:id", ArticleController.GetArticle)
		api.POST("/articles", ArticleController.AddArticle)
		api.DELETE("/articles/:id", ArticleController.DelArticle)
	}
}
