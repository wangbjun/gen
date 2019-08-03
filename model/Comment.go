package model

type Comment struct {
	Base
	UserID    uint
	ArticleId string
	Content   string
}
