package handler

import (
	"auth-service/internal/controller/model"
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
	apiGroup.GET("/user", userInfo)
	apiGroup.PUT("/refresh", refreshToken)
}

func signUp(ctx *gin.Context) {
	user := model.UserRequest{}

	if ctx.BindJSON(&user) != nil {
		return
	}

	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "miss password"})
		return
	}

	if !userService.IsValidEmail(user.Email) {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "email is invalid"})
		return
	}

	device := ctx.GetHeader("User-Agent")

	accessToken, refreshToken, err := userService.RegisterUser(user.Service, user.Email, user.Password, device)
	if errors.Is(err, service.ErrExistsEmail) {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	} else if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, model.TokenPairResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func signIn(ctx *gin.Context) {
	user := model.UserRequest{}

	if ctx.BindJSON(&user) != nil {
		return
	}

	device := ctx.GetHeader("User-Agent")

	accessToken, refreshToken, err := userService.Login(user.Service, user.Email, user.Password, device)
	if errors.Is(err, service.ErrIncorrectEmailOrPassword) {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	} else if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, model.TokenPairResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func userInfo(ctx *gin.Context) {
	defineUser := model.DefineUserRequest{}

	if ctx.BindJSON(&defineUser) != nil {
		return
	}

	email, expiresAt, err := userService.GetUserInfo(defineUser.Service, defineUser.AccessToken)
	if errors.Is(err, service.ErrTokenExpired) {
		ctx.JSON(http.StatusForbidden, model.ErrorResponse{Error: err.Error()})
		return
	} else if errors.Is(err, service.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error()})
		return
	} else if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, model.DefineUserResponse{Email: email, TokenExpiresAt: expiresAt})
}

func refreshToken(ctx *gin.Context) {
	tokenPair := model.TokenPairRequest{}

	if ctx.BindJSON(&tokenPair) != nil {
		return
	}

	accessToken, refreshToken, err := userService.RefreshTokens(tokenPair.AccessToken, tokenPair.RefreshToken)
	if errors.Is(err, service.ErrWrongRefreshToken) {
		ctx.JSON(http.StatusForbidden, model.ErrorResponse{Error: err.Error()})
		return
	} else if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, model.TokenPairResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
