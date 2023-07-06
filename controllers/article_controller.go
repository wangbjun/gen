package controllers

import (
	"gen/config"
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
		r.Failed(ctx, ParamError, translate(err))
		return
	}
	if article, err := r.ArticleService.Create(ctx, &param); err != nil {
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
		r.Failed(ctx, Failed, translate(err))
		return
	}
	param.Id, err = strconv.Atoi(ctx.Param("id"))
	if err != nil || param.Id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	if err := r.ArticleService.Update(ctx, &param); err != nil {
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
	article, err := r.ArticleService.GetById(ctx, id)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", article)
	}
	return
}

// GetAll 文章列表
func (r articleController) GetAll(ctx *gin.Context) {
	r.ParsePage(ctx)
	log.WithCtx(ctx).Info("GetAll Articles: " + config.Get().Env)
	articles, totalCount, err := r.ArticleService.GetAll(ctx, r.Page, r.PageSize)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", Pagination{
			List:       articles,
			Page:       r.Page,
			PageSize:   r.PageSize,
			TotalCount: totalCount,
		})
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
	err = r.ArticleService.Delete(ctx, id)
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
	err = r.ArticleService.AddComment(ctx, &param)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", "评论成功")
	}
	return
}
