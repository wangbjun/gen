package api

import (
	"gen/middleware"
)

func (hs *HTTPServer) registerRoutes() {
	hs.gin.GET("/", HttpServer.Index)
	r := hs.gin.Group("/api")
	{
		v1 := r.Group("/v1")
		{
			v1.Group("/articles").
				GET("", ArticleController.GetAll).     //文章列表
				GET("/:id", ArticleController.GetById) //文章详情

			v1.Group("/articles").Use(middleware.AuthMiddleware(hs.UserService)).
				POST("", ArticleController.Create).                 //添加文章
				PUT("/:id", ArticleController.Update).              //修改文章
				DELETE("/:id", ArticleController.Delete).           //删除文章
				POST("/:id/comments", ArticleController.AddComment) //添加评论

			v1.Group("/user").
				POST("/register", UserController.Register). //用户注册
				POST("/login", UserController.Login)        //用户登录

			v1.Group("/user").Use(middleware.AuthMiddleware(hs.UserService)).
				POST("/logout", UserController.Logout) //用户退出登录
		}
	}
}
