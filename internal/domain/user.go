package domain

import (
	"database/sql"
	"time"
)

type ID int64

type User struct {
	ID        ID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
