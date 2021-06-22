package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Id        uint `json:"id" gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	Salt      string
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserRegisterCommand struct {
	Name       string `form:"name" json:"name" binding:"gte=1,lte=20"`
	Email      string `form:"name" json:"email" binding:"required,email"`
	Password   string `form:"password" json:"password" binding:"required,gte=6"`
	RePassword string `form:"re_password" json:"re_password" binding:"eqfield=Password"`
}

type UserLoginCommand struct {
	Email    string `form:"name" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,gte=6"`
}

func IsUserEmailExisted(email string) (bool, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
