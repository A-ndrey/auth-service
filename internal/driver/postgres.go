package driver

import (
	"auth-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGorm() (*gorm.DB, error) {
	dsn := "user=auth-user password=auth-passwd dbname=auth_service host=127.0.0.1"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
