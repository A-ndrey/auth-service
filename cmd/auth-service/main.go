package main

import (
	"auth-service/internal/controller/handler"
	"auth-service/internal/driver"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	db, err := driver.NewPostgresGorm()
	if err != nil {
		log.Fatal(err)
	}

	userService := service.NewUserService(db)

	r := gin.Default()

	handler.RouteHealth(&r.RouterGroup)
	handler.RouteAPI(&r.RouterGroup, userService)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
