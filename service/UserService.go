package service

import (
	"errors"
	"fmt"
	"gen/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"time"
)

type UserError error

var (
	UserSecret               = []byte("@fc6951544^f55c644!@0d")
	UserExisted    UserError = errors.New("邮箱已存在")
	UserNotExisted UserError = errors.New("邮箱不存在")
	PasswordWrong  UserError = errors.New("邮箱或密码错误")
	LoginFailed    UserError = errors.New("登录失败")
)

type UserService interface {
	Register(name string, email string, password string) (string, error)
	Login(email string, password string) (string, error)
	ParseToken(string) (uint, error)
}

type userService struct {
	user *model.User
}

func NewUserService() UserService {
	return &userService{
		user: &model.User{},
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
	var user = model.User{}
	salt := GetUuidV4()[24:]
	user.Name = name
	user.Email = email
	user.Password = Sha1([]byte(password + salt))
	user.Salt = salt
	user.CreatedAt = uint(time.Now().Unix())
	user.UpdatedAt = uint(time.Now().Unix())

	err = model.DB.Save(&user).Error
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
	user, err := u.user.GetByEmail(email)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return "", UserNotExisted
		}
		return "", err
	}
	if user.Password != Sha1([]byte(password+user.Salt)) {
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
