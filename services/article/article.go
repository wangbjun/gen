package article

import (
	"errors"
	"gen/log"
	. "gen/models"
	"gen/registry"
	. "gen/services/sql_store"
	"github.com/jinzhu/gorm"
	"time"
)

type Service struct {
	SQLStore *SQLStore `inject:""`
}

func init() {
	registry.RegisterService(&Service{})
}

var (
	NotFound         = errors.New("文章不存在")
	PermissionDenied = errors.New("没有权限")
)

func (r Service) Init() error {
	return nil
}

func (r Service) Create(userId uint, form *ArticleCreateForm) (uint, error) {
	article := Article{
		Title:   form.Title,
		Content: form.Content,
		UserId:  userId,
	}
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	err := r.SQLStore.DB().Create(&article).Error
	if err != nil {
		return 0, err
	}
	return article.Id, nil
}

func (r Service) Edit(userId uint, param *ArticleUpdateForm) (uint, error) {
	var article Article
	err := r.SQLStore.DB().Where("id = ?", param.Id).Find(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, NotFound
	}
	if err != nil {
		return 0, err
	}
	if article.UserId != userId {
		return 0, PermissionDenied
	}
	article.Title = param.Title
	article.Content = param.Content
	article.UpdatedAt = time.Now()
	err = r.SQLStore.DB().Save(&article).Error
	if err != nil {
		return 0, err
	}
	return article.Id, nil
}

func (r Service) Detail(id uint) (*Article, error) {
	var article Article
	db := r.SQLStore.DB().Where("id = ?", id).Find(&article)
	if gorm.IsRecordNotFoundError(db.Error) {
		return nil, NotFound
	}
	if db.Error != nil {
		return nil, db.Error
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Logger.Sugar().Errorf("update view_num failed, error: %s", err)
			}
		}()
		db.UpdateColumn("view_num", gorm.Expr("view_num + 1"))
	}()
	return &article, nil
}

func (r Service) List(page int) ([]*Article, error) {
	var article []*Article
	offset := (page - 1) * 10
	err := r.SQLStore.DB().Limit(10).Offset(offset).
		Order("id desc").Find(&article).Error
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (r Service) Del(id int, userId uint) (bool, error) {
	var article Article
	err := r.SQLStore.DB().Where("id = ?", id).First(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, NotFound
	}
	if err != nil {
		return false, err
	}
	if article.UserId != userId {
		return false, PermissionDenied
	}
	err = r.SQLStore.DB().Delete(&article).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r Service) AddComment(userId uint, param *ArticleAddCommentForm) error {
	var article Article
	err := r.SQLStore.DB().Where("id = ?", param.ArticleId).First(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return NotFound
	}
	if err != nil {
		return err
	}
	comment := Comment{}
	comment.UserId = userId
	comment.ArticleId = param.ArticleId
	comment.Content = param.Content
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	err = r.SQLStore.DB().Create(&comment).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (r Service) ListComment(ArticleId uint) ([]*Comment, error) {
	var comments []*Comment
	err := r.SQLStore.DB().Where("article_id = ?", ArticleId).Find(&comments).Error
	if err != nil {
		return nil, err
	} else {
		return comments, nil
	}
}
