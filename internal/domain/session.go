package domain

import (
	"database/sql"
	"time"
)

type SessionID uint64

type Session struct {
	ID           SessionID
	UserID       UserID `gorm:"uniqueIndex:idx_sessions;uniqueIndex:idx_sessions_token"`
	Device       string `gorm:"uniqueIndex:idx_sessions"`
	RefreshToken string `gorm:"uniqueIndex:idx_sessions_token"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
