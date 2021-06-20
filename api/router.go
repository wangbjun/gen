package api

import (
	"gen/middleware"
)

func LoadRouter(hs *HTTPServer) {
	hs.Gin.GET("/", HttpServer.Index)

	r := hs.Gin.Group("/api")
	{
		v1 := r.Group("/v1")
		{
			v1.Group("/articles").
				GET("", ArticleController.ListArticle).   //文章列表
				GET("/:id", ArticleController.GetArticle) //文章详情

			v1.Group("/articles").Use(middleware.AuthMiddleware(hs.UserService)).
				POST("", ArticleController.CreateArticle).           //添加文章
				POST("/:id", ArticleController.EditArticle).         //修改文章
				DELETE("/:id", ArticleController.DelArticle).        //删除文章
				POST("/:id/comments", ArticleController.AddComment). //添加评论
				GET("/:id/comments", ArticleController.ListComment)  //评论列表

			v1.Group("/user").
				POST("/register", UserController.Register). //用户注册
				POST("/login", UserController.Login)        //用户登录

			v1.Group("/user").Use(middleware.AuthMiddleware(hs.UserService)).
				POST("/logout", UserController.Logout) //用户退出登录
		}
	}
}
