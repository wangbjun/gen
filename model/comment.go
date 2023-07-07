package model

import "time"

type Comment struct {
	Id        int        `json:"id" gorm:"primaryKey"`
	ArticleId int        `json:"article_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Comment) TableName() string {
	return "comments"
}
