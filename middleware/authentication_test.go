package middleware

import "testing"

func TextGetAuthMiddleware(t *testing.T) {
	t.Run("There is no Authorization header in the request", func(t *testing.T) {})
	t.Run(`The Authorization header doen't start with "Bearer "`, func(t *testing.T) {})
	t.Run(`Firebase app auth client is nil`, func(t *testing.T) {})
	t.Run(`Client.VerifyIDToken() returns an error`, func(t *testing.T) {})
	t.Run(`token has been set in gin.Context & Next() had been called`, func(t *testing.T) {})
}
