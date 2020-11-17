package main

import (
	"auth-service/internal/controller/handler"
	"auth-service/internal/driver"
	"auth-service/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
