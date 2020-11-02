package service

import (
	"auth-service/internal/domain"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net"
	"regexp"
	"strings"
)

type UserService interface {
	IsValidEmail(string) bool
	RegisterUser(string) error
}

type userService struct {
	gorm *gorm.DB
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewUserService(gorm *gorm.DB) UserService {
	return &userService{gorm: gorm}
}

func (u *userService) IsValidEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	if !emailRegex.MatchString(email) {
		return false
	}
	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

func (u *userService) RegisterUser(email string) error {
	user := domain.User{Email: email}

	result := u.gorm.Where(&user).First(&domain.User{})
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("email already registered")
	}

	u.gorm.Create(&user)

	return nil
}
