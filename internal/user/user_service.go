package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/anuragdaksh7/url-short-go/config"
	"github.com/anuragdaksh7/url-short-go/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	timeout time.Duration
	DB      *gorm.DB
}

func NewService() Service {
	return &service{
		time.Duration(20) * time.Second,
		config.DB,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hash),
	}
	res := config.DB.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	resp := &CreateUserRes{
		Id: user.ID,
	}

	return resp, nil
}

func (s *service) LoginUser(c context.Context, req *LoginReq) (*LoginRes, error) {
	_config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	var user model.User
	config.DB.First(&user, "email = ?", req.Email)

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(_config.JwtSecret))
	if err != nil {
		return nil, err
	}

	return &LoginRes{
		Token: tokenString,
	}, nil
}
