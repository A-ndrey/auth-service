package handler

import (
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

var userService service.UserService

func RouteAPI(group *gin.RouterGroup, service service.UserService) {
	userService = service
	apiGroup := group.Group("/api/v1")
	apiGroup.POST("/user", newUser)
}

func newUser(ctx *gin.Context) {
	var user struct {
		Email string `json:"email"`
	}

	if ctx.ShouldBindJSON(&user) != nil {
		ctx.JSON(http.StatusBadRequest, ErrorMessage{"wrong json format"})
		return
	}

	if !userService.IsValidEmail(user.Email) {
		ctx.JSON(http.StatusBadRequest, ErrorMessage{"email is invalid"})
		return
	}

	if err := userService.RegisterUser(user.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorMessage{err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}
