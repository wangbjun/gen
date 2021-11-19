package router

import (
	"gen/controllers"
	"gen/middleware"
	"gen/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/", controllers.BaseController.Index)
	r := server.Group("/api")
	{
		v1 := r.Group("/v1")
		{
			articleCtrl := controllers.ArticleController
			v1.Group("/articles").
				GET("", articleCtrl.GetAll).     //文章列表
				GET("/:id", articleCtrl.GetById) //文章详情

			userService := services.NewUserService()
			v1.Group("/articles").Use(middleware.AuthMiddleware(userService)).
				POST("", articleCtrl.Create).                 //添加文章
				PUT("/:id", articleCtrl.Update).              //修改文章
				DELETE("/:id", articleCtrl.Delete).           //删除文章
				POST("/:id/comments", articleCtrl.AddComment) //添加评论

			userCtrl := controllers.UserController
			v1.Group("/user").
				POST("/register", userCtrl.Register). //用户注册
				POST("/login", userCtrl.Login)        //用户登录

			v1.Group("/user").Use(middleware.AuthMiddleware(userService)).
				POST("/logout", userCtrl.Logout) //用户退出登录
		}
	}
}
