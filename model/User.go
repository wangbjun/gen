package model

import "github.com/jinzhu/gorm"

type User struct {
	Base
	Name     string
	Email    string
	Password string
	Salt     string
}

func (u User) IsEmailExisted(email string) (bool, error) {
	user, err := u.GetByEmail(email)
	if user == nil {
		return false, err
	} else {
		return true, nil
	}
}

func (User) GetByEmail(email string) (*User, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	} else {
		return &user, nil
	}
}
