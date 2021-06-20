package api

import (
	"gen/api/trans"
	"gen/log"
	"gen/models"
	"gen/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type articleController struct {
	*HTTPServer
}

var ArticleController = &articleController{HttpServer}

// CreateArticle 添加文章
func (r *articleController) CreateArticle(ctx *gin.Context) {
	var param models.ArticleCreateForm
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			r.Failed(ctx, ParamError, trans.Translate(e))
		} else {
			r.Failed(ctx, Failed, "请求错误")
		}
		return
	}
	userId := ctx.GetUint("userId")
	if userId <= 0 {
		r.Failed(ctx, NotLogin, "用户未登录")
		return
	}
	if id, err := r.HTTPServer.ArticleService.Create(userId, &param); err != nil {
		r.Failed(ctx, Failed, "添加文章失败")
	} else {
		r.Success(ctx, "添加文章成功", gin.H{"id": id})
	}
	return
}

// EditArticle 修改文章
func (r articleController) EditArticle(ctx *gin.Context) {
	var param models.ArticleUpdateForm
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			r.Failed(ctx, ParamError, trans.Translate(e))
		} else {
			r.Failed(ctx, Failed, "请求错误")
		}
		return
	}
	userId := ctx.GetUint("userId")
	if userId <= 0 {
		r.Failed(ctx, NotLogin, "用户未登录")
		return
	}
	if id, err := r.HTTPServer.ArticleService.Edit(userId, &param); err != nil {
		log.Error("edit article failed, error: " + err.Error())
		r.Failed(ctx, Failed, "修改文章失败")
	} else {
		r.Success(ctx, "修改文章成功", gin.H{"id": id})
	}
	return
}

// GetArticle 文章详情
func (r articleController) GetArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	at, err := r.HTTPServer.ArticleService.Detail(uint(id))
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", models.ArticleResult{
			Id:        at.Id,
			Title:     at.Title,
			Content:   at.Content,
			UserID:    at.UserId,
			ViewNum:   at.ViewNum,
			CreatedAt: at.CreatedAt.Format(utils.TimeFormatYmdHis),
			UpdatedAt: at.CreatedAt.Format(utils.TimeFormatYmdHis),
		})
	}
	return
}

// ListArticle 文章列表
func (r articleController) ListArticle(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	articles, err := r.HTTPServer.ArticleService.List(page)
	if err != nil {
		r.Failed(ctx, Failed, "获取文章列表失败")
	} else {
		var result = make([]models.ArticleResult, 0)
		for _, at := range articles {
			result = append(result, models.ArticleResult{
				Id:        at.Id,
				Title:     at.Title,
				Content:   at.Content,
				UserID:    at.UserId,
				ViewNum:   at.ViewNum,
				CreatedAt: at.CreatedAt.Format(utils.TimeFormatYmdHis),
				UpdatedAt: at.CreatedAt.Format(utils.TimeFormatYmdHis),
			})
		}
		r.Success(ctx, "ok", result)
	}
	return
}

// DelArticle 删除文章
func (r articleController) DelArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	_, err = r.HTTPServer.ArticleService.Del(id, ctx.GetUint("userId"))
	if err != nil {
		r.Failed(ctx, Failed, "删除失败")
	} else {
		r.Success(ctx, "删除成功", gin.H{"id": id})
	}
	return
}

// AddComment 新增评论
func (r articleController) AddComment(ctx *gin.Context) {
	var param models.ArticleAddCommentForm
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			r.Failed(ctx, ParamError, trans.Translate(e))
		} else {
			r.Failed(ctx, Failed, "请求错误")
		}
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	param.ArticleId = uint(id)
	err = r.HTTPServer.ArticleService.AddComment(ctx.GetUint("userId"), &param)
	if err != nil {
		r.Failed(ctx, Failed, "评论失败")
	} else {
		r.Success(ctx, "ok", "评论成功")
	}
	return
}

// ListComment 评论列表
func (r articleController) ListComment(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || articleId <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	comments, err := r.HTTPServer.ArticleService.ListComment(uint(articleId))
	if err != nil {
		r.Success(ctx, "ok", []string{})
	} else {
		r.Success(ctx, "ok", comments)
	}
	return
}
