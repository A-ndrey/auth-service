package service

import (
	"auth-service/internal/domain"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTService interface {
	NewToken(userID domain.UserID) (string, error)
	GetUserIDAndExpiresAt(token string) (domain.UserID, int64, error)
}

type jwtService struct {
	secret string
}

type UserClaims struct {
	UserID domain.UserID
	jwt.StandardClaims
}

const tokenDuration = 30 * time.Minute

var (
	ErrTokenExpired = errors.New("token expired")
)

func NewJWTService(secret string) JWTService {
	return &jwtService{secret: secret}
}

func (j *jwtService) NewToken(userID domain.UserID) (string, error) {
	claims := UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) GetUserIDAndExpiresAt(tokenString string) (domain.UserID, int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
	})
	if e, ok := err.(*jwt.ValidationError); ok && e.Errors&jwt.ValidationErrorExpired != 0 {
		return 0, 0, ErrTokenExpired
	} else if err != nil {
		return 0, 0, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return 0, 0, errors.New("invalid token")
	}

	return claims.UserID, claims.ExpiresAt, nil
}
