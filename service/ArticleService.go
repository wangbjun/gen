package service

import (
	"gen/model"
	"github.com/jinzhu/gorm"
	"time"
)

type ArticleService interface {
	New(article *model.Article) (uint, error)
	Get(id uint) (*model.Article, error)
	List(page int) ([]*model.Article, error)
	Del(id uint) (bool, error)
}

type articleService struct{}

func NewArticleService() ArticleService {
	return &articleService{}
}

const PageSize = 5

func (articleService) New(article *model.Article) (uint, error) {
	article.CreatedAt = uint(time.Now().Unix())
	article.UpdatedAt = uint(time.Now().Unix())
	err := model.DB.Create(&article).Error
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}
func (articleService) Get(id uint) (*model.Article, error) {
	var article model.Article
	err := model.DB.Where("id = ?", id).Find(&article).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
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

func (articleService) Del(id uint) (bool, error) {
	err := model.DB.Table("articles").Where("id = ?", id).
		Updates(map[string]interface{}{"status": 1, "updated_at": time.Now().Unix()}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
