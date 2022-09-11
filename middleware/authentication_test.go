package middleware

import (
	"net/http"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthMiddleware(t *testing.T) {
	t.Run("There is no Authorization header in the request", func(t *testing.T) {
		firebaseAuthClientMock, ginContextMock, errorHandlerMock := CreateMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(ErrNoAuthorizationHeader,
			"", http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Times(0)
		ginContextMock.EXPECT().GetHeader(AUTHORIZATION).Return("")
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ginContextMock)
	})

	t.Run(`The Authorization header doen't start with "Bearer "`, func(t *testing.T) {
		firebaseAuthClientMock, ginContextMock, errorHandlerMock := CreateMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(ErrAuthorizationHeaderDoesntStartWithBearer,
			"", http.StatusUnauthorized)
		ginContextMock.EXPECT().GetHeader(AUTHORIZATION).Return("ewhgfwwhgerwhg9")
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Times(0)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ginContextMock)
	})

	t.Run(`Firebase app auth client is nil`, func(t *testing.T) {
		_, ginContextMock, errorHandlerMock := CreateMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(ErrAuthClientIsNil,
			"", http.StatusInternalServerError)
		authMiddleware := GetAuthMiddleware(nil, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ginContextMock)
	})

	t.Run(`Client.VerifyIDToken() returns an error`, func(t *testing.T) {
		firebaseAuthClientMock, ginContextMock, errorHandlerMock := CreateMocks(t)
		ha := "eyJhbGciOiJ"
		ginContextMock.EXPECT().GetHeader(AUTHORIZATION).Return(BEARER + ha)
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError,
			"", http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(),
			ha).
			Return(nil, common.ErrError)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ginContextMock)
	})

	t.Run(`token doesn't have uid`, func(t *testing.T) {
		firebaseAuthClientMock, ginContextMock, errorHandlerMock := CreateMocks(t)
		ha := "eyJhbGciOiJ"
		tokeN := new(auth.Token)
		ginContextMock.EXPECT().GetHeader(AUTHORIZATION).Return(BEARER + ha)
		errorHandlerMock.EXPECT().HandleAppError(ErrNoUID, "", http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(),
			ha).
			Return(tokeN, nil)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ginContextMock)
	})

	t.Run(`token has been set in gin.Context & Next() had been called`, func(t *testing.T) {
		firebaseAuthClientMock, ginContextMock, errorHandlerMock := CreateMocks(t)
		ha := "eyJhbGciOiJ"
		tokeN := &auth.Token{UID: "woefhweh"}
		ginContextMock.EXPECT().GetHeader(AUTHORIZATION).Return(BEARER + ha)
		ginContextMock.EXPECT().Set(AuthToken, tokeN)
		ginContextMock.EXPECT().Next()
		errorHandlerMock.EXPECT().HandleAppError(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(),
			ha).
			Return(tokeN, nil)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		assert.NotNil(t, authMiddleware)
		authMiddleware(ginContextMock)
	})
}

func CreateMocks(t *testing.T) (*common.MockAuthClient, *common.MockWebContext, *common.MockErrorHandler) {
	mockCtrl := gomock.NewController(t)
	return common.NewMockAuthClient(mockCtrl), common.NewMockWebContext(mockCtrl), common.NewMockErrorHandler(mockCtrl)
}
