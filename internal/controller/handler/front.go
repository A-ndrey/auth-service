package handler

import (
	"auth-service/internal/controller/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

const title = "Auth"

func RouteFront(group *gin.RouterGroup) {
	webGroup := group.Group("/web")
	webGroup.Static("/assets", "front/assets")
	webGroup.Use(middleware.ServiceDefiner)
	webGroup.GET("/signup", web(true))
	webGroup.GET("/signin", web(false))
}

func web(isSignUp bool) func(*gin.Context) {
	return func(c *gin.Context) {
		service := c.GetString(middleware.ServiceKey)
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title":    title,
				"service":  service,
				"isSignUp": isSignUp,
			},
		)
	}
}
