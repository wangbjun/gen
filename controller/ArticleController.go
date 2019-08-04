package controller

import (
	"gen/model"
	"gen/service"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"strconv"
	"time"
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
	logs.Debug("add article")
	var (
		title, _   = c.GetPostForm("title")
		content, _ = c.GetPostForm("content")
	)
	if !govalidator.StringLength(title, "1", "100") {
		a.failed(c, ParamsError, "标题长度1-100")
		return
	}
	if !govalidator.StringLength(content, "1", "2000") {
		a.failed(c, ParamsError, "内容长度1-2000")
		return
	}
	userId := a.getUserId(c)
	if userId == 0 {
		a.failed(c, NotLogin, "未登录")
		return
	}
	var article = &model.Article{
		Title:   title,
		Content: content,
		UserID:  userId,
	}
	if id, err := a.articleService.New(article); err != nil {
		logs.Errorf("add article failed, error: " + err.Error())
		a.failed(c, Failed, "添加文章失败")
	} else {
		logs.Debugf("add article success，id:%d", id)
		a.success(c, "添加文章成功", map[string]interface{}{"id": id})
	}
	return
}

// 修改文章
func (a articleController) EditArticle(c *gin.Context) {
	logs.Debug("edit article")
	var (
		id, _      = strconv.Atoi(c.Param("id"))
		title, _   = c.GetPostForm("title")
		content, _ = c.GetPostForm("content")
	)
	if !govalidator.StringLength(title, "1", "100") {
		a.failed(c, ParamsError, "标题长度1-100")
		return
	}
	if !govalidator.StringLength(content, "1", "2000") {
		a.failed(c, ParamsError, "内容长度1-2000")
		return
	}
	userId := a.getUserId(c)
	if userId == 0 {
		a.failed(c, NotLogin, "未登录")
		return
	}
	var article = &model.Article{
		Title:   title,
		Content: content,
		UserID:  userId,
	}
	if id, err := a.articleService.Edit(uint(id), article); err != nil {
		logs.Errorf("edit article failed, error: " + err.Error())
		if _, ok := err.(service.ArticleError); ok {
			a.failed(c, UnAuthorized, err.Error())
		} else {
			a.failed(c, Failed, "修改文章失败")
		}
	} else {
		logs.Debugf("edit article success，id:%d", id)
		a.success(c, "修改文章成功", map[string]interface{}{"id": id})
	}
	return
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
		logs.Errorf("get article failed，id:%d, error:%s", id, err.Error())
		if _, ok := err.(service.ArticleError); ok {
			a.failed(c, Failed, err.Error())
		} else {
			a.failed(c, Failed, "获取文章失败")
		}
	} else {
		logs.Debugf("get article success，id:%d", id)
		a.success(c, "ok", resultWrap{
			Article:   article,
			CreatedAt: time.Unix(int64(article.CreatedAt), 0).Format("2006-01-02 15:04:05"),
			UpdatedAt: time.Unix(int64(article.UpdatedAt), 0).Format("2006-01-02 15:04:05"),
		})
		a.articleService.AddView(uint(id))
	}
	return
}

// 文章列表
func (a articleController) ListArticle(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	articles, err := a.articleService.List(page)
	if err != nil {
		logs.Errorf("list article failed，error:%s", err.Error())
		a.failed(c, Failed, "获取文章列表失败")
	} else {
		var result []resultWrap
		for _, article := range articles {
			result = append(result, resultWrap{
				Article:   article,
				CreatedAt: time.Unix(int64(article.CreatedAt), 0).Format("2006-01-02 15:04:05"),
				UpdatedAt: time.Unix(int64(article.UpdatedAt), 0).Format("2006-01-02 15:04:05"),
			})
		}
		if len(result) != 0 {
			a.success(c, "ok", result)
		} else {
			// 解决列表为空时，data为null的问题
			a.success(c, "ok", []string{})
		}
	}
	return
}

// 删除文章
func (a articleController) DelArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		a.failed(c, ParamsError, "id不能为空")
		return
	}
	_, err = a.articleService.Del(uint(id), a.getUserId(c))
	if err != nil {
		logs.Errorf("del article failed，id:%d, error:%s", id, err.Error())
		if _, ok := err.(service.ArticleError); ok {
			a.failed(c, Failed, err.Error())
		} else {
			a.failed(c, Failed, "删除失败")
		}
	} else {
		logs.Debugf("del article success，id:%d", id)
		a.success(c, "删除成功", map[string]interface{}{"id": id})
	}
	return
}

// 新增评论
func (a articleController) AddComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		a.failed(c, ParamsError, "id不能为空")
		return
	}
	content, _ := c.GetPostForm("content")
	if !govalidator.StringLength(content, "1", "500") {
		a.failed(c, ParamsError, "评论长度1-500")
		return
	}
	userId := a.getUserId(c)
	if userId == 0 {
		a.failed(c, NotLogin, "未登录")
		return
	}
	comment, err := a.articleService.AddComment(uint(id), userId, content)
	if err != nil {
		logs.Errorf("del article failed，id:%d, error:%s", id, err.Error())
		if _, ok := err.(service.ArticleError); ok {
			a.failed(c, NotFound, err.Error())
		} else {
			a.failed(c, Failed, "评论失败")
		}
	} else {
		a.success(c, "ok", comment)
	}
	return
}

func (a articleController) StopTheWorld() {
	go service.StopAdd()
	logs.Infof("stop the world,waiting for 3s")
	time.Sleep(time.Second * 3)
}
