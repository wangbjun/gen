package controller

import (
	"gen/model"
	"gen/service/article"
	"gen/zlog"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type articleController struct {
	Controller
	articleService article.Service
}

type articleResult struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	ViewNum   uint   `json:"view_num"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var ArticleController = &articleController{
	articleService: article.New(),
}

// 添加文章
func (r articleController) AddArticle(ctx *gin.Context) {
	r.LogSugar(ctx).Debug("add new mArticle")
	var (
		title, _   = ctx.GetPostForm("title")
		content, _ = ctx.GetPostForm("content")
	)
	if !govalidator.StringLength(title, "1", "100") {
		r.Failed(ctx, ParamError, "标题长度1-100")
		return
	}
	if !govalidator.StringLength(content, "1", "2000") {
		r.Failed(ctx, ParamError, "内容长度1-2000")
		return
	}
	mArticle := model.Article{}
	mArticle.Title = title
	mArticle.Content = content
	mArticle.UserID = ctx.MustGet("userId").(uint)
	mArticle.CreatedAt = time.Now()
	mArticle.UpdatedAt = time.Now()

	if id, err := r.articleService.Add(&mArticle); err != nil {
		r.LogSugar(ctx).Errorf("add mArticle Failed, error: " + err.Error())
		r.Failed(ctx, Failed, "添加文章失败")
	} else {
		r.LogSugar(ctx).Debugf("add mArticle Success，id: %d", id)
		r.Success(ctx, "添加文章成功", gin.H{"id": id})
	}
	return
}

// 修改文章
func (r articleController) EditArticle(ctx *gin.Context) {
	r.LogSugar(ctx).Debug("edit article")
	var (
		id, _      = strconv.Atoi(ctx.Param("id"))
		title, _   = ctx.GetPostForm("title")
		content, _ = ctx.GetPostForm("content")
	)
	if !govalidator.StringLength(title, "1", "100") {
		r.Failed(ctx, ParamError, "标题长度1-100")
		return
	}
	if !govalidator.StringLength(content, "1", "2000") {
		r.Failed(ctx, ParamError, "内容长度1-2000")
		return
	}
	at := model.Article{}
	at.Title = title
	at.Content = content
	at.UserID = ctx.MustGet("userId").(uint)
	at.UpdatedAt = time.Now()
	if id, err := r.articleService.Edit(uint(id), &at); err != nil {
		zlog.WithContext(ctx).Error("edit article Failed, error: " + err.Error())
		if _, ok := err.(article.Error); ok {
			r.Failed(ctx, NotFound, err.Error())
		} else {
			r.Failed(ctx, Failed, "修改文章失败")
		}
	} else {
		r.LogSugar(ctx).Debugf("edit article Success，id: %d", id)
		r.Success(ctx, "修改文章成功", gin.H{"id": id})
	}
	return
}

// 文章详情
func (r articleController) GetArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	at, err := r.articleService.Detail(uint(id))
	if err != nil {
		r.LogSugar(ctx).Errorf("get at Failed，id: %d, error: %s", id, err.Error())
		if _, ok := err.(article.Error); ok {
			r.Failed(ctx, NotFound, err.Error())
		} else {
			r.Failed(ctx, Failed, "获取文章失败")
		}
	} else {
		r.LogSugar(ctx).Debugf("get at Success，id: %d", id)
		r.Success(ctx, "ok", articleResult{
			Id:        at.ID,
			Title:     at.Title,
			Content:   at.Content,
			UserID:    at.UserID,
			ViewNum:   at.ViewNum,
			CreatedAt: at.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: at.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

// 文章列表
func (r articleController) ListArticle(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	articles, err := r.articleService.List(page)
	if err != nil {
		r.LogSugar(ctx).Errorf("list article Failed，error: %s", err.Error())
		r.Failed(ctx, Failed, "获取文章列表失败")
	} else {
		var result = make([]articleResult, 0)
		for _, at := range articles {
			result = append(result, articleResult{
				Id:        at.ID,
				Title:     at.Title,
				Content:   at.Content,
				UserID:    at.UserID,
				ViewNum:   at.ViewNum,
				CreatedAt: at.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: at.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		r.Success(ctx, "ok", result)
	}
	return
}

// 删除文章
func (r articleController) DelArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	_, err = r.articleService.Del(uint(id), uint(ctx.GetInt("userId")))
	if err != nil {
		r.LogSugar(ctx).Errorf("del article Failed，id: %d, error: %s", id, err.Error())
		if _, ok := err.(article.Error); ok {
			r.Failed(ctx, Failed, err.Error())
		} else {
			r.Failed(ctx, Failed, "删除失败")
		}
	} else {
		r.LogSugar(ctx).Debugf("del article Success，id: %d", id)
		r.Success(ctx, "删除成功", gin.H{"id": id})
	}
	return
}

// 新增评论
func (r articleController) AddComment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	content, _ := ctx.GetPostForm("content")
	if !govalidator.StringLength(content, "1", "500") {
		r.Failed(ctx, ParamError, "评论长度1-500")
		return
	}
	comment, err := r.articleService.AddComment(uint(id), uint(ctx.GetInt("userId")), content)
	if err != nil {
		r.LogSugar(ctx).Errorf("del article Failed，id: %d, error: %s", id, err.Error())
		if _, ok := err.(article.Error); ok {
			r.Failed(ctx, NotFound, err.Error())
		} else {
			r.Failed(ctx, Failed, "评论失败")
		}
	} else {
		r.Success(ctx, "ok", comment)
	}
	return
}

// 评论列表
func (r articleController) ListComment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		r.Failed(ctx, ParamError, "id不能为空")
		return
	}
	comments, err := r.articleService.ListComment(uint(id))
	if err != nil {
		r.LogSugar(ctx).Errorf("get article list Failed，error: %s", err.Error())
		r.Success(ctx, "ok", []string{})
	} else {
		r.Success(ctx, "ok", comments)
	}
	return
}
