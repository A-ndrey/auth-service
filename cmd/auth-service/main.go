package main

import (
	"auth-service/internal/controller/handler"
	"auth-service/internal/driver"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	db, err := driver.NewPostgresGorm()
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("env JWT_SECRET not assigned")
	}

	jwtService := service.NewJWTService(jwtSecret)
	sessionService := service.NewSessionService(db)
	userService := service.NewUserService(db, jwtService, sessionService)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	healthService := service.NewHealthService(sqlDB)

	r := gin.Default()

	handler.RouteHealth(&r.RouterGroup, healthService)
	handler.RouteAPI(&r.RouterGroup, userService)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
