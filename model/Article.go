package model

import "github.com/jinzhu/gorm"

type Article struct {
	gorm.Model
	Title   string `json:"title" valid:"min=1,max=100"`
	Content string `json:"content" valid:"required"`
	UserID  uint   `json:"user_id" valid:"required"`
	ViewNum uint   `json:"view_num"`
}
