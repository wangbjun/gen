package models

import (
	"time"
)

type Article struct {
	Id        uint       `json:"id" gorm:"primary_key"`
	Title     string     `json:"title" binding:"min=1,max=100"`
	Content   string     `json:"content" binding:"required"`
	UserId    uint       `json:"user_id" binding:"required"`
	ViewNum   uint       `json:"view_num"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ArticleCreateForm struct {
	Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type ArticleUpdateForm struct {
	Id      uint   `form:"id" json:"id" binding:"required"`
	Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
	Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type ArticleAddCommentForm struct {
	ArticleId uint   `form:"id" json:"id"`
	Content   string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type ArticleResult struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	ViewNum   uint   `json:"view_num"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
