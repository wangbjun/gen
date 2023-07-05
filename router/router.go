package router

import (
	"gen/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/", controllers.BaseController.Index)

	v1 := server.Group("/v1")
	articleCtrl := controllers.ArticleController
	{
		v1.GET("/articles", articleCtrl.GetAll)        //所有文章
		v1.POST("/articles", articleCtrl.Create)       //添加文章
		v1.GET("/articles/:id", articleCtrl.GetById)   //文章详情
		v1.PUT("/articles/:id", articleCtrl.Update)    //修改文章
		v1.DELETE("/articles/:id", articleCtrl.Delete) //删除文章

		v1.GET("/articles/:id/comments", articleCtrl.AddComment)  //获取文章评论
		v1.POST("/articles/:id/comments", articleCtrl.AddComment) //给文章添加评论
	}
}
