package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
)

const AUTHORIZATION string = "Authorization"
const AuthToken string = "AuthClaims"
const BEARER string = "Bearer "

var ErrNoAuthorizationHeader error = errors.New("there is no Authorization header in the web request")
var ErrAuthorizationHeaderDoesntStartWithBearer error = errors.New(`the Authorization header in the web request doesn't start with "Bearer "`)
var ErrAuthClientIsNil error = errors.New("auth client is nil")
var ErrIdTokenVerificationFailed error = errors.New("id token verification faild")
var ErrNoUID error = errors.New("there is no UID in the token")

func GetAuthMiddleware(authClient common.AuthClient, errorHandler common.ErrorHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if authClient == nil {
			errorHandler.HandleAppError(ctx, ErrAuthClientIsNil, http.StatusInternalServerError)
		} else {
			authorizationHeader := ctx.GetHeader(AUTHORIZATION)
			if authorizationHeader == "" {
				errorHandler.HandleAppError(ctx, ErrNoAuthorizationHeader, http.StatusUnauthorized)
			} else if !strings.HasPrefix(authorizationHeader, BEARER) {
				errorHandler.HandleAppError(ctx, ErrAuthorizationHeaderDoesntStartWithBearer, http.StatusUnauthorized)
			} else {
				token := strings.Replace(authorizationHeader, BEARER, "", 1)
				authToken, err := authClient.VerifyIDToken(context.Background(), token)
				if err != nil {
					errorHandler.HandleAppError(ctx, err, http.StatusUnauthorized)
				} else {
					if authToken.UID == "" {
						errorHandler.HandleAppError(ctx, ErrNoUID, http.StatusUnauthorized)
					} else {
						ctx.Set(AuthToken, authToken)
					}
				}
			}
		}
	}
}
