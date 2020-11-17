package middleware

import (
	"auth-service/internal/controller/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ServiceKey = "service"

func ServiceDefiner(ctx *gin.Context) {
	serviceQuery, ok := ctx.GetQuery(ServiceKey)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "empty 'service' query"})
		return
	}
	ctx.Set(ServiceKey, serviceQuery)
}
