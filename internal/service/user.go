package service

import (
	"auth-service/internal/domain"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net"
	"regexp"
	"strings"
)

type UserService interface {
	IsValidEmail(email string) bool
	RegisterUser(service, email, password string) (string, error)
}

type userService struct {
	db         *gorm.DB
	jwtService JWTService
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var (
	ErrExistsEmail = errors.New("email already exists")
)

func NewUserService(db *gorm.DB, jwtService JWTService) UserService {
	return &userService{db: db, jwtService: jwtService}
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

func (u *userService) RegisterUser(service, email, password string) (string, error) {
	user := domain.User{Service: service, Email: email}

	result := u.db.Where(&user).First(&domain.User{})
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", ErrExistsEmail
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	user.HashedPassword = hashedPassword

	token, err := u.jwtService.NewToken(service, email)
	if err != nil {
		return "", err
	}

	result = u.db.Create(&user)
	if result.Error != nil {
		return "", result.Error
	}

	return token, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}