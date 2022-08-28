package middleware

import "github.com/gin-gonic/gin"

func GetAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
