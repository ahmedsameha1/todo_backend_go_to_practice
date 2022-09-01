package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
)

var AuthToken string = "AuthClaims"
var ErrNoAuthorizationHeader error = errors.New("there is no Authorization header in the web request")
var ErrAuthorizationHeaderDoesntStartWithBearer error = errors.New(`the Authorization header in the web request doesn't start with "Bearer "`)
var ErrAuthClientIsNil error = errors.New("auth client is nil")
var ErrIdTokenVerificationFailed error = errors.New("id token verification faild")
var ErrError error = errors.New("an error")

func GetAuthMiddleware(authClient common.AuthClient, errorHandler common.ErrorHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if authClient == nil {
			errorHandler.HandleAppError(ErrAuthClientIsNil, "", http.StatusInternalServerError)
		} else {
			authorizationHeader := ctx.Request.Header.Get("Authorization")
			if authorizationHeader == "" {
				errorHandler.HandleAppError(ErrNoAuthorizationHeader, "", http.StatusUnauthorized)
			} else if !strings.HasPrefix(authorizationHeader, "Bearer ") {
				errorHandler.HandleAppError(ErrAuthorizationHeaderDoesntStartWithBearer, "", http.StatusUnauthorized)
			} else {
				token := strings.Replace(authorizationHeader, "Bearer ", "", 1)
				authToken, err := authClient.VerifyIDToken(context.Background(), token)
				if err != nil {
					errorHandler.HandleAppError(err, "", http.StatusUnauthorized)
				} else {
					ctx.Set(AuthToken, authToken)
				}
			}
		}
	}
}
