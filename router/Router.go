package router

import (
	. "gen/controller"
	"gen/middleware"
	"github.com/gin-gonic/gin"
)

func Route(Router *gin.Engine) {
	Router.Use(middleware.Request()) //全局中间件,记录请求日志

	Router.GET("/", BaseController.Index)

	Router.Group("/api")

	Router.Group("/api/v1/articles").
		GET("/", ArticleController.ListArticle).  //文章列表
		GET("/:id", ArticleController.GetArticle) //文章详情

	Router.Group("/api/v1/articles").Use(middleware.Auth()).
		POST("/", ArticleController.AddArticle).             //添加文章
		POST("/:id", ArticleController.EditArticle).         //修改文章
		DELETE("/:id", ArticleController.DelArticle).        //删除文章
		POST("/:id/comments", ArticleController.AddComment). //添加评论
		GET("/:id/comments", ArticleController.ListComment)  //评论列表

	Router.Group("/api/v1/user").
		POST("/register", UserController.Register). //用户注册
		POST("/login", UserController.Login)        //用户登录

	Router.Group("/api/v1/user").Use(middleware.Auth()).
		POST("/logout", UserController.Logout) //用户退出登录

}
