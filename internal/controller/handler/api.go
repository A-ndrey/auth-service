package handler

import (
	"auth-service/internal/controller/presenter"
	"auth-service/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var userService service.UserService

func RouteAPI(group *gin.RouterGroup, service service.UserService) {
	userService = service
	apiGroup := group.Group("/api/v1")
	apiGroup.POST("/signup", signUp)
	apiGroup.POST("/signin", signIn)
}

func signUp(ctx *gin.Context) {
	user := presenter.UserRequest{}

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
	if errors.Is(err, service.ErrExistsEmail) {
		ctx.JSON(http.StatusBadRequest, presenter.ErrorResponse{Error: err.Error()})
		return
	} else if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, presenter.TokenPairResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func signIn(ctx *gin.Context) {
	user := presenter.UserRequest{}

	if ctx.ShouldBindJSON(&user) != nil {
		ctx.JSON(http.StatusBadRequest, presenter.ErrorResponse{Error: "wrong json format"})
		return
	}

	device := ctx.GetHeader("User-Agent")

	accessToken, refreshToken, err := userService.Login(user.Service, user.Email, user.Password, device)
	if errors.Is(err, service.ErrIncorrectEmailOrPassword) {
		ctx.JSON(http.StatusBadRequest, presenter.ErrorResponse{Error: err.Error()})
	} else if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, presenter.TokenPairResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
