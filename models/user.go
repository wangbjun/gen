package models

import (
	"time"
)

type User struct {
	Id        uint `json:"id" gorm:"primary_key"`
	Name      string
	Email     string
	Password  string
	Salt      string
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserRegisterForm struct {
	Name       string `form:"name" json:"name" binding:"gte=1,lte=20"`
	Email      string `form:"name" json:"email" binding:"required,email"`
	Password   string `form:"password" json:"password" binding:"required,gte=6"`
	RePassword string `form:"re_password" json:"re_password" binding:"eqfield=Password"`
}

type UserLoginForm struct {
	Email    string `form:"name" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,gte=6"`
}
