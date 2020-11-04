package domain

import (
	"database/sql"
	"time"
)

type ID uint64

type User struct {
	ID             ID
	Service        string `gorm:"uniqueIndex:idx_user"`
	Email          string `gorm:"uniqueIndex:idx_user"`
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
}
