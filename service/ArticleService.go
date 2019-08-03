package service

import (
	"gen/model"
	"time"
)

type ArticleService interface {
	New(article *model.Article) (uint, error)
	Get(id uint) (*model.Article, error)
	List(page int) ([]*model.Article, error)
	Del(id uint) (bool, error)
}

type service struct{}

func NewArticleService() ArticleService {
	return &service{}
}

const PageSize = 5

func (service) New(article *model.Article) (uint, error) {
	article.CreatedAt = uint(time.Now().Unix())
	article.UpdatedAt = uint(time.Now().Unix())
	err := model.DB.Create(&article).Error
	if err != nil {
		return 0, model.DB.Error
	}
	return article.ID, nil
}
func (service) Get(id uint) (*model.Article, error) {
	var article model.Article
	err := model.DB.First(&article).Where("id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (service) List(page int) ([]*model.Article, error) {
	var article []*model.Article
	err := model.DB.Where("status = 0").Limit(PageSize).
		Offset(page).Order("id desc").
		Find(&article).Error
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (service) Del(id uint) (bool, error) {
	err := model.DB.Table("articles").Where("id = ?", id).
		Updates(map[string]interface{}{"status": 1, "updated_at": time.Now().Unix()}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
