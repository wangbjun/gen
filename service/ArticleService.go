package service

import (
	"errors"
	"gen/model"
	"github.com/jinzhu/gorm"
	logs "github.com/sirupsen/logrus"
	"time"
)

type ArticleService interface {
	New(*model.Article) (uint, error)
	Edit(uint, *model.Article) (uint, error)
	Get(id uint) (*model.Article, error)
	List(page int) ([]*model.Article, error)
	Del(id uint, userId uint) (bool, error)
	AddView(id uint)
	AddComment(id uint, userId uint, content string) (*model.Comment, error)
	ListComment(id uint) ([]*model.Comment, error)
}

var (
	addView = make(chan uint)
	addCtrl = make(chan string)
)

type articleService struct{}

func NewArticleService() ArticleService {
	return &articleService{}
}

func init() {
	saveTicker()
}

type ArticleError error

var (
	ArticleNotFound  ArticleError = errors.New("文章不存在")
	PermissionDenied ArticleError = errors.New("没有权限")
)

func (a articleService) AddView(id uint) {
	go func() {
		addView <- id
	}()
}

func StopAdd() {
	go func() {
		addCtrl <- "shutdown"
	}()
}

func (articleService) New(article *model.Article) (uint, error) {
	article.CreatedAt = uint(time.Now().Unix())
	article.UpdatedAt = uint(time.Now().Unix())
	err := model.DB.Create(&article).Error
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (a articleService) Edit(id uint, params *model.Article) (uint, error) {
	var article model.Article
	err := model.DB.Where("id = ?", id).Find(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, ArticleNotFound
	}
	if article.UserID != params.UserID {
		return 0, PermissionDenied
	}
	article.Title = params.Title
	article.Content = params.Content
	article.UpdatedAt = uint(time.Now().Unix())
	err = model.DB.Save(&article).Error
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (articleService) Get(id uint) (*model.Article, error) {
	var article model.Article
	err := model.DB.Where("id = ?", id).Find(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, ArticleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (articleService) List(page int) ([]*model.Article, error) {
	var article []*model.Article
	offset := (page - 1) * PageSize
	err := model.DB.Where("status = 0").Limit(PageSize).
		Offset(offset).Order("id desc").
		Find(&article).Error
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (articleService) Del(id uint, userId uint) (bool, error) {
	var article model.Article
	err := model.DB.Where("id = ?", id).First(&article).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return false, ArticleNotFound
	}
	if article.UserID != userId {
		return false, PermissionDenied
	}
	article.Status = StatusDel
	article.UpdatedAt = uint(time.Now().Unix())
	err = model.DB.Save(&article).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a articleService) AddComment(id uint, userId uint, content string) (*model.Comment, error) {
	var article model.Article
	err := model.DB.Where("id = ?", id).First(&article).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil, ArticleNotFound
	}
	comment := model.Comment{}
	comment.UserID = userId
	comment.ArticleId = id
	comment.Content = content
	comment.CreatedAt = uint(time.Now().Unix())
	comment.UpdatedAt = uint(time.Now().Unix())
	err = model.DB.Create(&comment).Error
	if err != nil {
		return nil, err
	} else {
		return &comment, nil
	}
}

func (a articleService) ListComment(id uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := model.DB.Where("article_id = ? and status = 0", id).Find(&comments).Error
	if err != nil {
		return nil, err
	} else {
		return comments, nil
	}
}

// 保存浏览数
func saveTicker() {
	go func() {
		var times = make(map[uint]int)
		for {
			select {
			case id := <-addView:
				times[id]++
				if len(times) > 1000 {
					go saveView(times)
					times = make(map[uint]int)
				}
			case c := <-addCtrl:
				if c == "save" && len(times) > 0 {
					logs.Infof("begin save view times")
					go saveView(times)
					times = make(map[uint]int)
				} else if c == "shutdown" && len(times) > 0 {
					logs.Infof("begin save view times before shutdown")
					go saveView(times)
				}
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(5 * 60 * time.Second)
		for {
			<-ticker.C
			addCtrl <- "save"
		}
	}()
}
func saveView(times map[uint]int) {
	var article model.Article
	for id, count := range times {
		model.DB.Model(&article).Where("id = ?", id).UpdateColumn("view_num", gorm.Expr("view_num + ?", count))
	}
}
