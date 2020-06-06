package article

import (
	"errors"
	"gen/lib/zlog"
	. "gen/model"
	"github.com/jinzhu/gorm"
	"time"
)

type Service interface {
	Add(*Article) (uint, error)
	Edit(uint, *Article) (uint, error)
	Detail(id uint) (*Article, error)
	List(page int) ([]*Article, error)
	Del(id uint, userId uint) (bool, error)
	AddComment(id uint, userId uint, content string) (*Comment, error)
	ListComment(id uint) ([]*Comment, error)
}

type articleService struct{}

func New() Service {
	return &articleService{}
}

type Error struct {
	error
}

func articleError(err string) Error {
	return Error{errors.New(err)}
}

var (
	NotFound         = articleError("文章不存在")
	PermissionDenied = articleError("没有权限")
)

func (articleService) Add(article *Article) (uint, error) {
	err := DB().Create(&article).Error
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (articleService) Edit(id uint, params *Article) (uint, error) {
	var article Article
	err := DB().Where("id = ?", id).Find(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, NotFound
	}
	if err != nil {
		return 0, err
	}
	if article.UserID != params.UserID {
		return 0, PermissionDenied
	}
	article.Title = params.Title
	article.Content = params.Content
	article.UpdatedAt = time.Now()
	err = DB().Save(&article).Error
	if err != nil {
		return 0, err
	}
	return article.ID, nil
}

func (articleService) Detail(id uint) (*Article, error) {
	var article Article
	db := DB().Where("id = ?", id).Find(&article)
	if gorm.IsRecordNotFoundError(db.Error) {
		return nil, NotFound
	}
	if db.Error != nil {
		return nil, db.Error
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zlog.Logger.Sugar().Errorf("update view_num failed, error: %s", err)
			}
		}()
		db.UpdateColumn("view_num", gorm.Expr("view_num + 1"))
	}()
	return &article, nil
}

func (articleService) List(page int) ([]*Article, error) {
	var article []*Article
	offset := (page - 1) * 10
	err := DB().Limit(10).Offset(offset).Order("id desc").Find(&article).Error
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (articleService) Del(id uint, userId uint) (bool, error) {
	var article Article
	err := DB().Where("id = ?", id).First(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, NotFound
	}
	if err != nil {
		return false, err
	}
	if article.UserID != userId {
		return false, PermissionDenied
	}
	err = DB().Delete(&article).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a articleService) AddComment(id uint, userId uint, content string) (*Comment, error) {
	var article Article
	err := DB().Where("id = ?", id).First(&article).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, NotFound
	}
	if err != nil {
		return nil, nil
	}
	comment := Comment{}
	comment.UserID = userId
	comment.ArticleId = id
	comment.Content = content
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	err = DB().Create(&comment).Error
	if err != nil {
		return nil, err
	} else {
		return &comment, nil
	}
}

func (a articleService) ListComment(id uint) ([]*Comment, error) {
	var comments []*Comment
	err := DB().Where("article_id = ? and status = 0", id).Find(&comments).Error
	if err != nil {
		return nil, err
	} else {
		return comments, nil
	}
}
