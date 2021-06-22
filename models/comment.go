package models

import "time"

type Comment struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	UserId    int        `json:"user_id"`
	ArticleId int        `json:"article_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
