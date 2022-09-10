package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
)

const AUTHORIZATION string = "Authorization"
const AuthToken string = "AuthClaims"
const BEARER string = "Bearer "

var ErrNoAuthorizationHeader error = errors.New("there is no Authorization header in the web request")
var ErrAuthorizationHeaderDoesntStartWithBearer error = errors.New(`the Authorization header in the web request doesn't start with "Bearer "`)
var ErrAuthClientIsNil error = errors.New("auth client is nil")
var ErrIdTokenVerificationFailed error = errors.New("id token verification faild")
var ErrNoUIDinToken error = errors.New("there is no UID in the token")

func GetAuthMiddleware(authClient common.AuthClient, errorHandler common.ErrorHandler) func(common.WebContext) {
	return func(ctx common.WebContext) {
		if authClient == nil {
			errorHandler.HandleAppError(ErrAuthClientIsNil, "", http.StatusInternalServerError)
		} else {
			authorizationHeader := ctx.GetHeader(AUTHORIZATION)
			if authorizationHeader == "" {
				errorHandler.HandleAppError(ErrNoAuthorizationHeader, "", http.StatusUnauthorized)
			} else if !strings.HasPrefix(authorizationHeader, BEARER) {
				errorHandler.HandleAppError(ErrAuthorizationHeaderDoesntStartWithBearer, "", http.StatusUnauthorized)
			} else {
				token := strings.Replace(authorizationHeader, BEARER, "", 1)
				authToken, err := authClient.VerifyIDToken(context.Background(), token)
				if err != nil {
					errorHandler.HandleAppError(err, "", http.StatusUnauthorized)
				} else {
					if authToken.UID == "" {
						errorHandler.HandleAppError(ErrNoUIDinToken, "", http.StatusUnauthorized)
					} else {
						ctx.Set(AuthToken, authToken)
						ctx.Next()
					}
				}
			}
		}
	}
}
