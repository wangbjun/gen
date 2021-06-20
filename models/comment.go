package models

import "time"

type Comment struct {
	Id        uint       `json:"id" gorm:"primary_key"`
	UserId    uint       `json:"user_id"`
	ArticleId uint       `json:"article_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
