package controller

import (
	"gen/model"
	"gen/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
	"unicode/utf8"
)

type articleController struct {
	Controller
	articleService service.ArticleService
}

type resultWrap struct {
	*model.Article
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var ArticleController = &articleController{
	articleService: service.NewArticleService(),
}

// 添加文章
func (a articleController) AddArticle(c *gin.Context) {
	log.Debug("add article")
	var (
		title, _   = c.GetPostForm("title")
		content, _ = c.GetPostForm("content")
	)
	if utf8.RuneCountInString(title) < 1 || utf8.RuneCountInString(title) > 100 {
		a.failed(c, ParamsError, "标题长度1-100")
		return
	}
	if utf8.RuneCountInString(content) < 1 || utf8.RuneCountInString(content) > 2000 {
		a.failed(c, ParamsError, "内容长度1-2000")
		return
	}
	var article = &model.Article{
		Title:   title,
		Content: content,
		UserID:  uint(c.GetInt("userId")),
	}
	if id, err := a.articleService.New(article); err != nil {
		log.Error("add article failed, error: " + err.Error())
		a.failed(c, Failed, "添加文章失败")
		return
	} else {
		log.Debugf("add article success，id:%d", id)
		a.success(c, "添加文章成功", map[string]interface{}{"id": id})
		return
	}
}

// 文章详情
func (a articleController) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		a.failed(c, ParamsError, "id不能为空")
		return
	}
	article, err := a.articleService.Get(uint(id))
	if err != nil {
		log.Errorf("get article failed，id:%d, error:%s", id, err.Error())
		a.failed(c, Failed, "获取文章失败")
		return
	}
	if article.ID == 0 {
		a.failed(c, NotFound, "文章不存在")
		return
	} else {
		log.Debugf("get article success，id:%d", id)
		a.success(c, "ok", resultWrap{
			Article:   article,
			CreatedAt: time.Unix(int64(article.CreatedAt), 0).Format("2006-01-02 15:04:05"),
			UpdatedAt: time.Unix(int64(article.UpdatedAt), 0).Format("2006-01-02 15:04:05"),
		})
	}
}

// 文章列表
func (a articleController) ListArticle(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	articles, err := a.articleService.List(page)
	if err != nil {
		log.Errorf("list article failed，error:%s", err.Error())
		a.failed(c, Failed, "获取文章列表失败")
		return
	} else {
		var result []resultWrap
		for _, article := range articles {
			result = append(result, resultWrap{
				Article:   article,
				CreatedAt: time.Unix(int64(article.CreatedAt), 0).Format("2006-01-02 15:04:05"),
				UpdatedAt: time.Unix(int64(article.UpdatedAt), 0).Format("2006-01-02 15:04:05"),
			})
		}
		a.success(c, "ok", result)
		return
	}
}

// 删除文章
func (a articleController) DelArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		a.failed(c, ParamsError, "id不能为空")
		return
	}
	_, err = a.articleService.Del(uint(id))
	if err != nil {
		log.Errorf("del article failed，id:%d, error:%s", id, err.Error())
		a.failed(c, Failed, "删除失败")
		return
	} else {
		log.Debugf("del article success，id:%d", id)
		a.success(c, "删除成功", map[string]interface{}{"id": id})
	}
}
