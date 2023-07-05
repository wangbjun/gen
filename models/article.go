package models

import (
	"context"
	"gorm.io/gorm"
	"time"
)

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

func CreateArticle(ctx context.Context, param *CreateArticleCommand) (*Article, error) {
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

func DeleteArticle(ctx context.Context, id int) error {
	article := Article{
		Id: id,
	}
	err := NewOrm(ctx).Delete(&article).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateArticle(ctx context.Context, param *UpdateArticleCommand) error {
	article := Article{
		Id:        param.Id,
		Title:     param.Title,
		Content:   param.Content,
		UpdatedAt: time.Now(),
	}
	err := NewOrm(ctx).Updates(&article).Error
	if err != nil {
		return err
	}
	return nil
}

func GetArticleById(ctx context.Context, id int) (*Article, error) {
	var article Article
	err := NewOrm(ctx).Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func GetArticles(ctx context.Context, page, pageSize int) ([]*Article, int64, error) {
	var totalCount int64
	var articles []*Article
	err := NewOrm(ctx).Model(Article{}).Count(&totalCount).
		Limit(pageSize).Offset((page - 1) * pageSize).Order("id desc").Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}
	return articles, totalCount, nil
}

func CreateArticleComment(ctx context.Context, param *CreateArticleCommentCommand) error {
	comment := Comment{}
	comment.ArticleId = param.Id
	comment.Content = param.Content
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	err := NewOrm(ctx).Create(&comment).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func AddViewNum(ctx context.Context, id int) error {
	err := NewOrm(ctx).Model(Article{}).Where("id = ?", id).
		UpdateColumn("view_num", gorm.Expr("view_num + 1")).Error
	if err != nil {
		return err
	}
	return nil
}
