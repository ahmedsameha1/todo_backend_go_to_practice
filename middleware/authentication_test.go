package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthMiddleware(t *testing.T) {
	t.Run("There is no Authorization header in the request", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		authMiddleware := GetAuthMiddleware()
		assert.NotNil(t, authMiddleware)
	})
	t.Run(`The Authorization header doen't start with "Bearer "`, func(t *testing.T) {})
	t.Run(`Firebase app auth client is nil`, func(t *testing.T) {})
	t.Run(`Client.VerifyIDToken() returns an error`, func(t *testing.T) {})
	t.Run(`token has been set in gin.Context & Next() had been called`, func(t *testing.T) {})
}
