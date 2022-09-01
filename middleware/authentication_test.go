package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthMiddleware(t *testing.T) {
	t.Run("There is no Authorization header in the request", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		r := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(r)
		req, _ := http.NewRequest("GET", "/todos", nil)
		ctx.Request = req
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		firebaseAuthClientMock := common.NewMockAuthClient(mockCtrl)
		errorHandlerMock.EXPECT().HandleAppError(ErrNoAuthorizationHeader,
			"", http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Times(0)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ctx)
	})

	t.Run(`The Authorization header doen't start with "Bearer "`, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		r := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(r)
		req, _ := http.NewRequest("GET", "/todos", nil)
		ctx.Request = req
		ctx.Request.Header.Add("Authorization", "etewetweew")
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		firebaseAuthClientMock := common.NewMockAuthClient(mockCtrl)
		errorHandlerMock.EXPECT().HandleAppError(ErrAuthorizationHeaderDoesntStartWithBearer,
			"", http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Times(0)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ctx)
	})

	t.Run(`Firebase app auth client is nil`, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		r := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(r)
		req, _ := http.NewRequest("GET", "/todos", nil)
		ctx.Request = req
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		var firebaseAuthClient common.AuthClient = nil
		errorHandlerMock.EXPECT().HandleAppError(ErrAuthClientIsNil,
			"", http.StatusInternalServerError)
		authMiddleware := GetAuthMiddleware(firebaseAuthClient, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ctx)
	})

	t.Run(`Client.VerifyIDToken() returns an error`, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		r := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(r)
		req, _ := http.NewRequest("GET", "/todos", nil)
		ctx.Request = req
		ctx.Request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5C")
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		firebaseAuthClientMock := common.NewMockAuthClient(mockCtrl)
		errorHandlerMock.EXPECT().HandleAppError(ErrError,
			"", http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(),
		 "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5C").
		 Return(nil, ErrError)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ctx)
	})

	t.Run(`token has been set in gin.Context & Next() had been called`, func(t *testing.T) {})
}
