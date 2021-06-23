package models

import (
	"gorm.io/gorm"
	"time"
)

func CreateArticle(param *CreateArticleCommand) (*Article, error) {
	article := Article{
		Title:   param.Title,
		Content: param.Content,
		UserId:  param.UserId,
	}
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	err := db.Create(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func DeleteArticle(id int) error {
	article := Article{
		Id: id,
	}
	err := db.Delete(&article).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateArticle(param *UpdateArticleCommand) error {
	article := Article{
		Id:        param.Id,
		Title:     param.Title,
		Content:   param.Content,
		UpdatedAt: time.Now(),
	}
	err := db.Updates(&article).Error
	if err != nil {
		return err
	}
	return nil
}

func GetArticleById(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func GetArticles(page, pageSize int) ([]*Article, error) {
	var articles []*Article
	err := db.Limit(pageSize).Offset((page - 1) * pageSize).
		Order("id desc").Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func CreateArticleComment(param *CreateArticleCommentCommand) error {
	comment := Comment{}
	comment.UserId = param.UserId
	comment.ArticleId = param.Id
	comment.Content = param.Content
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	err := db.Create(&comment).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func AddViewNum(id int) error {
	err := db.Model(&Article{}).Where("id = ?", id).
		UpdateColumn("view_num", gorm.Expr("view_num + 1")).Error
	if err != nil {
		return err
	}
	return nil
}

type Article struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	Title     string     `json:"title" binding:"min=1,max=100"`
	Content   string     `json:"content" binding:"required"`
	UserId    int        `json:"user_id" binding:"required"`
	ViewNum   int        `json:"view_num"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateArticleCommand struct {
	Id      int
	UserId  int
	Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type UpdateArticleCommand struct {
	Id      int
	UserId  int
	Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type CreateArticleCommentCommand struct {
	UserId  int
	Id      int    `form:"id" json:"id"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}
