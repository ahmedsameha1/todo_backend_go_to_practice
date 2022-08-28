package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthMiddleware(t *testing.T) {
	t.Run("There is no Authorization header in the request", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		authMiddleware := GetAuthMiddleware()
		assert.NotNil(t, authMiddleware)
		r := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(r)
		req, _ := http.NewRequest("GET", "/todos", nil)
		ctx.Request = req
		authMiddleware(ctx)
	})
	t.Run(`The Authorization header doen't start with "Bearer "`, func(t *testing.T) {})
	t.Run(`Firebase app auth client is nil`, func(t *testing.T) {})
	t.Run(`Client.VerifyIDToken() returns an error`, func(t *testing.T) {})
	t.Run(`token has been set in gin.Context & Next() had been called`, func(t *testing.T) {})
}
