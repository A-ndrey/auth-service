package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTService interface {
	NewToken(service, email string) (string, error)
	CheckToken(token string, service, email string) error
}

type jwtService struct {
	secret string
}

type UserClaims struct {
	Service string
	Email   string
	jwt.StandardClaims
}

const tokenDuration = 30 * time.Minute

func NewJWTService(secret string) JWTService {
	return &jwtService{secret: secret}
}

func (j *jwtService) NewToken(service, email string) (string, error) {
	claims := UserClaims{
		Service: service,
		Email:   email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) CheckToken(tokenString string, service, email string) error {
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

	if claims.Service != service || claims.Email != email {
		return errors.New("token doesn't correspond user")
	}

	return nil
}
