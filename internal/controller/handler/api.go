package handler

import (
	"auth-service/internal/controller/presenter"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var userService service.UserService

func RouteAPI(group *gin.RouterGroup, service service.UserService) {
	userService = service
	apiGroup := group.Group("/api/v1")
	apiGroup.POST("/user", newUser)
}

func newUser(ctx *gin.Context) {
	var user struct {
		Service  string `json:"service"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if ctx.ShouldBindJSON(&user) != nil {
		ctx.JSON(http.StatusBadRequest, presenter.ErrorResponse{Error: "wrong json format"})
		return
	}

	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, presenter.ErrorResponse{Error: "miss password"})
		return
	}

	if !userService.IsValidEmail(user.Email) {
		ctx.JSON(http.StatusBadRequest, presenter.ErrorResponse{Error: "email is invalid"})
		return
	}

	device := ctx.GetHeader("User-Agent")

	accessToken, refreshToken, err := userService.RegisterUser(user.Service, user.Email, user.Password, device)
	if err == service.ErrExistsEmail {
		ctx.JSON(http.StatusBadRequest, presenter.ErrorResponse{Error: service.ErrExistsEmail.Error()})
		return
	} else if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, presenter.NewUserResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
