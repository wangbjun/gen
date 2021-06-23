package user

import (
	"errors"
	"fmt"
	. "gen/models"
	"gen/registry"
	"gen/utils"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"time"
)

var (
	Secret        = []byte("@fc6951544^f55c644!@0d")
	Existed       = errors.New("邮箱已存在")
	NotExisted    = errors.New("邮箱不存在")
	PasswordWrong = errors.New("邮箱或密码错误")
	LoginFailed   = errors.New("登录失败")
)

type UserService struct {
	SQLStore *SQLService `inject:""`
}

func init() {
	registry.RegisterService(&UserService{})
}

func (r UserService) Init() error {
	return nil
}

func (r UserService) Register(param *UserRegisterCommand) (string, error) {
	emailExisted, err := IsUserEmailExisted(param.Email)
	if err != nil {
		return "", err
	}
	if emailExisted {
		return "", Existed
	}
	var user = User{}
	salt := utils.GetUuidV4()[24:]
	user.Name = param.Name
	user.Email = param.Email
	user.Password = utils.Sha1([]byte(param.Password + salt))
	user.Salt = salt
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = DB().Save(&user).Error
	if err != nil {
		return "", err
	} else {
		token, err := r.createToken(user.Id)
		if err != nil {
			return "", err
		} else {
			return token, nil
		}
	}
}

func (r UserService) Login(param *UserLoginCommand) (string, error) {
	var user User
	err := DB().Where("email = ?", param.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", NotExisted
		}
		return "", err
	}
	if user.Password != utils.Sha1([]byte(param.Password+user.Salt)) {
		return "", PasswordWrong
	} else {
		token, err := r.createToken(user.Id)
		if err != nil {
			return "", LoginFailed
		} else {
			return token, nil
		}
	}
}

// ParseToken 解析token
func (r UserService) ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return Secret, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["userId"].(float64)), nil
	} else {
		return 0, err
	}
}

// 创建token
func (r UserService) createToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
