package sqlstore

import (
	"gen/models"
	"github.com/jinzhu/gorm"
)

func IsUserEmailExisted(email string) (bool, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
