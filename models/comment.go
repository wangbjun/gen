package models

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	ArticleId uint   `json:"article_id"`
	Content   string `json:"content"`
}
