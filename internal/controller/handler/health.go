package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RouteHealth(group *gin.RouterGroup) {
	healthGroup := group.Group("/health")
	ping(healthGroup)
}

func ping(group *gin.RouterGroup) {
	group.GET("/ping", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
}
