package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if authorizationHeader := ctx.Request.Header.Get("Authorization"); authorizationHeader == "" {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	}
}
