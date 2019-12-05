package controller

import (
	"gen/library/zlog"
	"gen/model"
	"gen/service/articleService"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type articleController struct {
	Controller
	articleService articleService.Service
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
	articleService: articleService.New(),
}

// 添加文章
func (ac articleController) AddArticle(ctx *gin.Context) {
	zlog.WithContext(ctx).Sugar().Debug("add article")
	var (
		title, _   = ctx.GetPostForm("title")
		content, _ = ctx.GetPostForm("content")
	)
	if !govalidator.StringLength(title, "1", "100") {
		ac.failed(ctx, ParamError, "标题长度1-100")
		return
	}
	if !govalidator.StringLength(content, "1", "2000") {
		ac.failed(ctx, ParamError, "内容长度1-2000")
		return
	}
	article := model.Article{}
	article.Title = title
	article.Content = content
	article.UserID = uint(ctx.GetInt("userId"))
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	if id, err := ac.articleService.Add(&article); err != nil {
		zlog.WithContext(ctx).Sugar().Errorf("add article failed, error: " + err.Error())
		ac.failed(ctx, Failed, "添加文章失败")
	} else {
		zlog.WithContext(ctx).Sugar().Debugf("add article success，id: %d", id)
		ac.success(ctx, "添加文章成功", gin.H{"id": id})
	}
	return
}

// 修改文章
func (ac articleController) EditArticle(ctx *gin.Context) {
	zlog.WithContext(ctx).Sugar().Debug("edit article")
	var (
		id, _      = strconv.Atoi(ctx.Param("id"))
		title, _   = ctx.GetPostForm("title")
		content, _ = ctx.GetPostForm("content")
	)
	if !govalidator.StringLength(title, "1", "100") {
		ac.failed(ctx, ParamError, "标题长度1-100")
		return
	}
	if !govalidator.StringLength(content, "1", "2000") {
		ac.failed(ctx, ParamError, "内容长度1-2000")
		return
	}
	article := model.Article{}
	article.Title = title
	article.Content = content
	article.UserID = uint(ctx.GetInt("userId"))
	article.UpdatedAt = time.Now()
	if id, err := ac.articleService.Edit(uint(id), &article); err != nil {
		zlog.WithContext(ctx).Error("edit article failed, error: " + err.Error())
		if _, ok := err.(articleService.Error); ok {
			ac.failed(ctx, NotFound, err.Error())
		} else {
			ac.failed(ctx, Failed, "修改文章失败")
		}
	} else {
		zlog.WithContext(ctx).Sugar().Debugf("edit article success，id: %d", id)
		ac.success(ctx, "修改文章成功", map[string]interface{}{"id": id})
	}
	return
}

// 文章详情
func (ac articleController) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		ac.failed(c, ParamError, "id不能为空")
		return
	}
	article, err := ac.articleService.Detail(uint(id))
	if err != nil {
		zlog.WithContext(c).Sugar().Errorf("get article failed，id: %d, error: %s", id, err.Error())
		if _, ok := err.(articleService.Error); ok {
			ac.failed(c, NotFound, err.Error())
		} else {
			ac.failed(c, Failed, "获取文章失败")
		}
	} else {
		zlog.WithContext(c).Sugar().Debugf("get article success，id: %d", id)
		ac.success(c, "ok", articleResult{
			Id:        article.ID,
			Title:     article.Title,
			Content:   article.Content,
			UserID:    article.UserID,
			ViewNum:   article.ViewNum,
			CreatedAt: article.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: article.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return
}

// 文章列表
func (ac articleController) ListArticle(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	articles, err := ac.articleService.List(page)
	if err != nil {
		zlog.WithContext(c).Sugar().Errorf("list article failed，error: %s", err.Error())
		ac.failed(c, Failed, "获取文章列表失败")
	} else {
		var result []articleResult
		for _, article := range articles {
			result = append(result, articleResult{
				Id:        article.ID,
				Title:     article.Title,
				Content:   article.Content,
				UserID:    article.UserID,
				ViewNum:   article.ViewNum,
				CreatedAt: article.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: article.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		if len(result) != 0 {
			ac.success(c, "ok", result)
		} else {
			// 解决列表为空时，data为null的问题
			ac.success(c, "ok", []string{})
		}
	}
	return
}

// 删除文章
func (ac articleController) DelArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ac.failed(ctx, ParamError, "id不能为空")
		return
	}
	_, err = ac.articleService.Del(uint(id), uint(ctx.GetInt("userId")))
	if err != nil {
		zlog.WithContext(ctx).Sugar().Errorf("del article failed，id: %d, error: %s", id, err.Error())
		if _, ok := err.(articleService.Error); ok {
			ac.failed(ctx, Failed, err.Error())
		} else {
			ac.failed(ctx, Failed, "删除失败")
		}
	} else {
		zlog.WithContext(ctx).Sugar().Debugf("del article success，id: %d", id)
		ac.success(ctx, "删除成功", map[string]interface{}{"id": id})
	}
	return
}

// 新增评论
func (ac articleController) AddComment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ac.failed(ctx, ParamError, "id不能为空")
		return
	}
	content, _ := ctx.GetPostForm("content")
	if !govalidator.StringLength(content, "1", "500") {
		ac.failed(ctx, ParamError, "评论长度1-500")
		return
	}
	comment, err := ac.articleService.AddComment(uint(id), uint(ctx.GetInt("userId")), content)
	if err != nil {
		zlog.WithContext(ctx).Sugar().Errorf("del article failed，id: %d, error: %s", id, err.Error())
		if _, ok := err.(articleService.Error); ok {
			ac.failed(ctx, NotFound, err.Error())
		} else {
			ac.failed(ctx, Failed, "评论失败")
		}
	} else {
		ac.success(ctx, "ok", comment)
	}
	return
}

// 评论列表
func (ac articleController) ListComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		ac.failed(c, ParamError, "id不能为空")
		return
	}
	comments, err := ac.articleService.ListComment(uint(id))
	if err != nil {
		zlog.WithContext(c).Sugar().Errorf("get article list failed，error: %s", err.Error())
		ac.success(c, "ok", []string{})
	} else {
		ac.success(c, "ok", comments)
	}
	return
}
