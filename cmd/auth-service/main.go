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
	userService := service.NewUserService(db, jwtService)

	r := gin.Default()

	handler.RouteHealth(&r.RouterGroup)
	handler.RouteAPI(&r.RouterGroup, userService)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
