package handler

import (
	"auth-service/internal/controller/model"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RouteHealth(group *gin.RouterGroup, service service.HealthService) {
	healthGroup := group.Group("/health")
	healthGroup.GET("/", status(service))
	healthGroup.GET("/ping", ping)
}

func status(healthService service.HealthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		health := model.HealthResponse{
			DBConnection: healthService.DBConnectionStatus(),
			WorkingTime:  healthService.WorkingTime().String(),
		}

		ctx.JSON(http.StatusOK, health)
	}
}

func ping(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
