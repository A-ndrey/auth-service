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
	RegisterUser(service, email, password, device string) (string, string, error)
	Login(service, email, password, device string) (string, string, error)
	GetUserInfo(service, token string) (string, int64, error)
	RefreshTokens(accessToken, refreshToken string) (string, string, error)
}

type userService struct {
	db             *gorm.DB
	jwtService     JWTService
	sessionService SessionService
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var (
	ErrExistsEmail              = errors.New("email already exists")
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrUserNotFound             = errors.New("user not found")
)

func NewUserService(db *gorm.DB, jwtService JWTService, sessionService SessionService) UserService {
	return &userService{db: db, jwtService: jwtService, sessionService: sessionService}
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

func (u *userService) RegisterUser(service, email, password, device string) (string, string, error) {
	user := domain.User{Service: service, Email: email}

	result := u.db.First(&domain.User{}, &user)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", "", ErrExistsEmail
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", "", err
	}

	user.HashedPassword = hashedPassword

	var accessToken, refreshToken string

	err = u.db.Transaction(func(tx *gorm.DB) error {
		result = u.db.Create(&user)
		if result.Error != nil {
			return result.Error
		}

		accessToken, err = u.jwtService.NewToken(user.ID)
		if err != nil {
			return err
		}

		refreshToken, err = u.sessionService.NewSession(user.ID, device)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userService) Login(service, email, password, device string) (string, string, error) {
	user := domain.User{Service: service, Email: email}

	result := u.db.First(&user, &user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", "", ErrIncorrectEmailOrPassword
	} else if result.Error != nil {
		return "", "", result.Error
	}

	if !IsCorrectPassword(password, user.HashedPassword) {
		return "", "", ErrIncorrectEmailOrPassword
	}

	accessToken, err := u.jwtService.NewToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.sessionService.NewSession(user.ID, device)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userService) GetUserInfo(service, token string) (string, int64, error) {
	userID, expiresAt, err := u.jwtService.GetUserIDAndExpiresAt(token)
	if err != nil {
		return "", 0, err
	}

	user := domain.User{ID: userID, Service: service}
	result := u.db.Find(&user, &user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", 0, ErrUserNotFound
	} else if result.Error != nil {
		return "", 0, result.Error
	}

	return user.Email, expiresAt, nil
}

func (u *userService) RefreshTokens(accessToken, refreshToken string) (string, string, error) {
	userID, _, err := u.jwtService.GetUserIDAndExpiresAt(accessToken)
	if !errors.Is(err, ErrTokenExpired) && err != nil {
		return "", "", err
	}

	newAccessToken, err := u.jwtService.NewToken(userID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := u.sessionService.CheckAndUpdateToken(userID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func IsCorrectPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
