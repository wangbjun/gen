package model

type Comment struct {
	Base
	UserID    uint   `json:"user_id"`
	ArticleId uint   `json:"article_id"`
	Content   string `json:"content"`
}
