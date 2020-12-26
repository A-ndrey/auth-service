package driver

import (
	"auth-service/internal/domain"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

const schemaName = "auth_service"

func NewPostgresGorm() (*gorm.DB, error) {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return nil, errors.New("env POSTGRES_USER not assigned")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return nil, errors.New("env POSTGRES_PASSWORD not assigned")
	}

	dbname, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return nil, errors.New("env POSTGRES_DB not assigned")
	}

	dbHost, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return nil, errors.New("env POSTGRES_HOST not assigned")
	}

	config := gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: schemaName + "."},
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s", user, password, dbname, dbHost)
	db, err := gorm.Open(postgres.Open(dsn), &config)
	if err != nil {
		return nil, err
	}

	result := db.Exec("create schema if not exists " + schemaName)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Session{}); err != nil {
		return nil, err
	}

	return db, nil
}
