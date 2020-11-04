package service

import (
	"auth-service/internal/domain"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

type SessionService interface {
	GetActiveSessions(userID domain.UserID) ([]domain.Session, error)
	NewSession(userID domain.UserID, device string) (string, error)
	CheckAndUpdateToken(userID domain.UserID, token string) (string, error)
	DeleteSession(userID domain.UserID, device string) error
}

type sessionService struct {
	db *gorm.DB
}

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~=+%^*/()[]{}/!@#$?|"
	length  = 20
)

func NewSessionService(db *gorm.DB) SessionService {
	return &sessionService{db: db}
}

func (s *sessionService) GetActiveSessions(userID domain.UserID) ([]domain.Session, error) {
	var sessions []domain.Session
	result := s.db.Find(&sessions, &domain.Session{UserID: userID})
	if result.Error != nil {
		return nil, result.Error
	}

	return sessions, nil
}

func (s *sessionService) NewSession(userID domain.UserID, device string) (string, error) {
	session := domain.Session{
		UserID: userID,
		Device: device,
	}

	result := s.db.Find(&session, &session)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", result.Error
	}

	session.RefreshToken = GenerateToken()

	result = s.db.Save(&session)
	if result.Error != nil {
		return "", result.Error
	}

	return session.RefreshToken, nil
}

func (s *sessionService) CheckAndUpdateToken(userID domain.UserID, token string) (string, error) {
	session := domain.Session{
		UserID:       userID,
		RefreshToken: token,
	}

	result := s.db.Find(&session, &session)
	if result.Error != nil {
		return "", result.Error
	}

	session.RefreshToken = GenerateToken()

	result = s.db.Save(&session)
	if result.Error != nil {
		return "", result.Error
	}

	return session.RefreshToken, nil
}

func (s *sessionService) DeleteSession(userID domain.UserID, device string) error {
	session := domain.Session{
		UserID: userID,
		Device: device,
	}

	result := s.db.Delete(&session)

	return result.Error
}

func GenerateToken() string {
	rand.Seed(time.Now().Unix())
	var builder strings.Builder
	runeCharset := []rune(charset)
	for i := 0; i < length; i++ {
		builder.WriteRune(runeCharset[rand.Intn(len(runeCharset))])
	}

	return builder.String()
}
