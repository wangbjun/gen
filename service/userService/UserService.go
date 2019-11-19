package userService

import (
	"errors"
	"fmt"
	"gen/common"
	. "gen/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"time"
)

type UserError struct {
	error
}

func userError(err string) UserError {
	return UserError{errors.New(err)}
}

var (
	UserSecret     = []byte("@fc6951544^f55c644!@0d")
	UserExisted    = userError("邮箱已存在")
	UserNotExisted = userError("邮箱不存在")
	PasswordWrong  = userError("邮箱或密码错误")
	LoginFailed    = userError("登录失败")
)

type Service interface {
	Register(name string, email string, password string) (string, error)
	Login(email string, password string) (string, error)
	ParseToken(string) (uint, error)
}

type userService struct {
	user *User
}

func New() Service {
	return &userService{
		user: &User{},
	}
}

func (u userService) Register(name string, email string, password string) (string, error) {
	emailExisted, err := u.user.IsEmailExisted(email)
	if err != nil {
		return "", err
	}
	if emailExisted {
		return "", UserExisted
	}
	var user = User{}
	salt := common.GetUuidV4()[24:]
	user.Name = name
	user.Email = email
	user.Password = common.Sha1([]byte(password + salt))
	user.Salt = salt
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = DB.Save(&user).Error
	if err != nil {
		return "", err
	} else {
		token, err := u.createToken(user.ID)
		if err != nil {
			return "", err
		} else {
			return token, nil
		}
	}
}

func (u userService) Login(email string, password string) (string, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return "", UserNotExisted
		}
		return "", err
	}
	if user.Password != common.Sha1([]byte(password+user.Salt)) {
		return "", PasswordWrong
	} else {
		token, err := u.createToken(user.ID)
		if err != nil {
			return "", LoginFailed
		} else {
			return token, nil
		}
	}
}

// 解析token
func (u userService) ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return UserSecret, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["userId"].(float64)), nil
	} else {
		return 0, err
	}
}

// 创建token
func (u userService) createToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(UserSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
