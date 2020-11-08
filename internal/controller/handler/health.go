package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RouteHealth(group *gin.RouterGroup) {
	healthGroup := group.Group("/health")
	healthGroup.GET("/ping", ping)
}

func ping(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
