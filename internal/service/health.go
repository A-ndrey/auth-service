package service

import (
	"database/sql"
	"time"
)

type HealthService interface {
	WorkingTime() time.Duration
	DBConnectionStatus() string
}

type healthService struct {
	db        *sql.DB
	startedAt time.Time
}

func NewHealthService(db *sql.DB) HealthService {
	return &healthService{db: db, startedAt: time.Now()}
}

func (h *healthService) WorkingTime() time.Duration {
	return time.Since(h.startedAt)
}

func (h *healthService) DBConnectionStatus() string {
	if err := h.db.Ping(); err != nil {
		return err.Error()
	}

	return "OK"
}
