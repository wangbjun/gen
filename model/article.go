package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ArticleModel = Article{}

type Article struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	Title     string     `json:"title" binding:"min=1,max=100"`
	Content   string     `json:"content" binding:"required"`
	ViewNum   int        `json:"view_num"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Article) TableName() string {
	return "articles"
}

type CreateArticleCommand struct {
	Id      int
	Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type UpdateArticleCommand struct {
	Id      int
	Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type CreateArticleCommentCommand struct {
	Id      int    `form:"id" json:"id"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

func (Article) GetAll(ctx context.Context, page, pageSize int) ([]*Article, int64, error) {
	var totalCount int64
	var articles []*Article
	err := NewOrm(ctx).Model(Article{}).Count(&totalCount).
		Limit(pageSize).Offset((page - 1) * pageSize).Order("id desc").Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}
	return articles, totalCount, nil
}

func (Article) GetById(ctx context.Context, id int) (*Article, error) {
	var article Article
	err := NewOrm(ctx).Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (Article) Create(ctx context.Context, param *CreateArticleCommand) (*Article, error) {
	article := Article{
		Title:     param.Title,
		Content:   param.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := NewOrm(ctx).Create(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (Article) Delete(ctx context.Context, id int) error {
	article := Article{
		Id: id,
	}
	return NewOrm(ctx).Delete(&article).Error
}

func (Article) Update(ctx context.Context, param *UpdateArticleCommand) error {
	article := Article{
		Id:        param.Id,
		Title:     param.Title,
		Content:   param.Content,
		UpdatedAt: time.Now(),
	}
	return NewOrm(ctx).Updates(&article).Error
}

func (Article) CreateComment(ctx context.Context, param *CreateArticleCommentCommand) error {
	comment := Comment{
		ArticleId: param.Id,
		Content:   param.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return NewOrm(ctx).Create(&comment).Error
}

func (Article) AddViewNum(ctx context.Context, id int) error {
	return NewOrm(ctx).Model(Article{}).Where("id = ?", id).
		UpdateColumn("view_num", gorm.Expr("view_num + 1")).Error
}
