package controllers

import (
	"gen/log"
	"gen/models"
	"gen/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type articleController struct {
	*Controller
	*services.ArticleService
}

var ArticleController = articleController{
	Controller:     BaseController,
	ArticleService: services.NewArticleService(),
}

// Create 添加文章
func (r articleController) Create(ctx *gin.Context) {
	var param models.CreateArticleCommand
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		r.Failed(ctx, ParamError, err.Error())
		return
	}
	param.UserId = ctx.GetInt("userId")
	if param.UserId <= 0 {
		r.Failed(ctx, NotLogin, "用户未登录")
		return
	}
	if article, err := r.ArticleService.Create(&param); err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "添加文章成功", article)
	}
	return
}

// Update 修改文章
func (r articleController) Update(ctx *gin.Context) {
	var param models.UpdateArticleCommand
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		r.Failed(ctx, Failed, "请求错误")
		return
	}
	param.Id, err = strconv.Atoi(ctx.Param("id"))
	if err != nil || param.Id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	param.UserId = ctx.GetInt("userId")
	if param.UserId <= 0 {
		r.Failed(ctx, NotLogin, "用户未登录")
		return
	}
	if err := r.ArticleService.Update(&param); err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "修改文章成功", nil)
	}
	return
}

// GetById 文章详情
func (r articleController) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	article, err := r.ArticleService.GetById(id)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", article)
	}
	return
}

// GetAll 文章列表
func (r articleController) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	log.WithCtx(ctx).Info("GetAll Articles")
	articles, err := r.ArticleService.GetAll(ctx, page)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", articles)
	}
	return
}

// Delete 删除文章
func (r articleController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	err = r.ArticleService.Delete(id, ctx.GetInt("userId"))
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "删除成功", nil)
	}
	return
}

// AddComment 新增评论
func (r articleController) AddComment(ctx *gin.Context) {
	var param models.CreateArticleCommentCommand
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		r.Failed(ctx, Failed, "请求错误")
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	param.Id = id
	err = r.ArticleService.AddComment(&param)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", "评论成功")
	}
	return
}
