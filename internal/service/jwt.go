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
	CheckToken(token string, userID domain.UserID) error
}

type jwtService struct {
	secret string
}

type UserClaims struct {
	UserID domain.UserID
	jwt.StandardClaims
}

const tokenDuration = 30 * time.Minute

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

func (j *jwtService) CheckToken(tokenString string, userID domain.UserID) error {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.secret, nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token")
	}

	if claims.UserID != userID {
		return errors.New("token doesn't correspond userID")
	}

	return nil
}
