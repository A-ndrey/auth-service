package handler

import (
	"auth-service/internal/controller/middleware"
	"auth-service/internal/controller/model"
	"auth-service/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/lane-c-wagner/go-password-validator"
	"log"
	"net/http"
	"strings"
)

func RouteAPI(group *gin.RouterGroup, service service.UserService) {
	apiGroup := group.Group("/api/v1")

	apiGroup.POST("/password/check", checkPassword)
	apiGroup.PUT("/refresh", refreshToken(service))

	srvcDefGroup := apiGroup.Use(middleware.ServiceDefiner)

	srvcDefGroup.POST("/signup", signUp(service))
	srvcDefGroup.POST("/signin", signIn(service))
	srvcDefGroup.GET("/user", userInfo(service))
}

func signUp(userService service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		serviceQuery := ctx.GetString(middleware.ServiceKey)
		device := ctx.GetHeader("User-Agent")

		accessToken, refreshToken, err := userService.RegisterUser(serviceQuery, user.Email, user.Password, device)
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
}

func signIn(userService service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := model.UserRequest{}

		if ctx.BindJSON(&user) != nil {
			return
		}

		serviceQuery := ctx.GetString(middleware.ServiceKey)
		device := ctx.GetHeader("User-Agent")

		accessToken, refreshToken, err := userService.Login(serviceQuery, user.Email, user.Password, device)
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
}

func userInfo(userService service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		accessToken, err := getToken(authHeader)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
			return
		}

		serviceQuery := ctx.GetString(middleware.ServiceKey)

		email, expiresAt, err := userService.GetUserInfo(serviceQuery, accessToken)
		if errors.Is(err, service.ErrTokenExpired) {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
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
}

func refreshToken(userService service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshRequest := model.RefreshRequest{}

		if ctx.BindJSON(&refreshRequest) != nil {
			return
		}

		accessToken, refreshToken, err := userService.RefreshTokens(refreshRequest.RefreshToken)
		if errors.Is(err, service.ErrWrongRefreshToken) {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
			return
		} else if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, model.TokenPairResponse{AccessToken: accessToken, RefreshToken: refreshToken})
	}
}

func checkPassword(ctx *gin.Context) {
	request := model.PasswordCheckRequest{}

	if ctx.BindJSON(&request) != nil {
		return
	}

	response := model.PasswordCheckResponse{MinStrength: model.PasswordVeryWeak, MaxStrength: model.PasswordVeryStrong}

	entropy := passwordvalidator.GetEntropy(request.Password)
	if entropy < 40 {
		response.Strength = model.PasswordVeryWeak
	} else if entropy < 50 {
		response.Strength = model.PasswordWeak
	} else if entropy < 60 {
		response.Strength = model.PasswordMedium
	} else if entropy < 70 {
		response.Strength = model.PasswordStrong
	} else {
		response.Strength = model.PasswordVeryStrong
	}

	err := passwordvalidator.Validate(request.Password, 60)
	if err != nil {
		response.Recommendation = err.Error()
	}

	ctx.JSON(http.StatusOK, response)
}

func getToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("no 'Authorization' header")
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("no 'Bearer' prefix")
	}

	return authHeader[len(prefix):], nil
}
