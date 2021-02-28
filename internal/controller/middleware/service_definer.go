package middleware

import (
	"auth-service/internal/controller/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ServiceKey = "service"

func ServiceDefiner(ctx *gin.Context) {
	serviceQuery, _ := ctx.GetQuery(ServiceKey)
	if serviceQuery == "" {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			model.ErrorResponse{Error: "empty 'service' query", ErrorType: model.CommonError},
		)
		return
	}
	ctx.Set(ServiceKey, serviceQuery)
}
