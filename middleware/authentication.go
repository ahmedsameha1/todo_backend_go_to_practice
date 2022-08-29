package middleware

import (
	"errors"
	"net/http"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
)
var ErrNoAuthorizationHeader error = errors.New("there is no Authorization header in the web request")

func GetAuthMiddleware(errorHandler common.ErrorHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if authorizationHeader := ctx.Request.Header.Get("Authorization"); authorizationHeader == "" {
			errorHandler.HandleAppError(ErrNoAuthorizationHeader, "", http.StatusUnauthorized)
		}
	}
}
