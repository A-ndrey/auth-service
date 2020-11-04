package domain

import (
	"database/sql"
	"time"
)

type UserID uint64

type User struct {
	ID             UserID
	Service        string `gorm:"uniqueIndex:idx_users"`
	Email          string `gorm:"uniqueIndex:idx_users"`
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
}
