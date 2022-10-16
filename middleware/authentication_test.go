package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthMiddleware(t *testing.T) {
	t.Run("There is no Authorization header in the request", func(t *testing.T) {
		firebaseAuthClientMock, gin_context, errorHandlerMock := CreateMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, ErrNoAuthorizationHeader,
			http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Times(0)
		web_request := &http.Request{
			Header: map[string][]string{AUTHORIZATION: {""}}}
		gin_context.Request = web_request
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		authMiddleware(gin_context)
	})

	t.Run(`The Authorization header doen't start with "Bearer "`, func(t *testing.T) {
		firebaseAuthClientMock, gin_context, errorHandlerMock := CreateMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, ErrAuthorizationHeaderDoesntStartWithBearer,
			http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(), gomock.Any()).Times(0)
		web_request := &http.Request{
			Header: map[string][]string{AUTHORIZATION: {"ewhgfwwhgerwhg9"}}}
		gin_context.Request = web_request
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		authMiddleware(gin_context)
	})

	t.Run(`Firebase app auth client is nil`, func(t *testing.T) {
		_, gin_context, errorHandlerMock := CreateMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, ErrAuthClientIsNil,
			http.StatusInternalServerError)
		authMiddleware := GetAuthMiddleware(nil, errorHandlerMock)
		authMiddleware(gin_context)
	})

	t.Run(`Client.VerifyIDToken() returns an error`, func(t *testing.T) {
		firebaseAuthClientMock, gin_context, errorHandlerMock := CreateMocks(t)
		ha := "eyJhbGciOiJ"
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError,
			http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(),
			ha).
			Return(nil, common.ErrError)
		web_request := &http.Request{
			Header: map[string][]string{AUTHORIZATION: {BEARER + ha}}}
		gin_context.Request = web_request
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		authMiddleware(gin_context)
	})

	t.Run(`token doesn't have uid`, func(t *testing.T) {
		firebaseAuthClientMock, gin_context, errorHandlerMock := CreateMocks(t)
		ha := "eyJhbGciOiJ"
		tokeN := new(auth.Token)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, ErrNoUID, http.StatusUnauthorized)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(),
			ha).
			Return(tokeN, nil)
		web_request := &http.Request{
			Header: map[string][]string{AUTHORIZATION: {BEARER + ha}}}
		gin_context.Request = web_request
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		authMiddleware(gin_context)
	})

	t.Run(`token has been set in gin.Context & the next middleware had been called`, func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		_, gin_engine := gin.CreateTestContext(httptest.NewRecorder())
		firebaseAuthClientMock, _, errorHandlerMock := CreateMocks(t)
		ha := "eyJhbGciOiJ"
		tokeN := &auth.Token{UID: "woefhweh"}
		errorHandlerMock.EXPECT().HandleAppError(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		firebaseAuthClientMock.EXPECT().VerifyIDToken(gomock.Any(),
			ha).
			Return(tokeN, nil)
		authMiddleware := GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
		gin_engine.Use(authMiddleware, func(ctx *gin.Context) { ctx.Set("next", "called") })
		gin_engine.GET("/", func(ctx *gin.Context) {
			v, _ := ctx.Get("next")
			n, _ := ctx.Get(AuthToken)
			ctx.JSON(http.StatusOK, gin.H{"token": n, "next": v})
		})
		server := httptest.NewServer(gin_engine)
		defer server.Close()
		req, err := http.NewRequest("GET", "http://"+server.Listener.Addr().String()+"/", nil)
		if err != nil {
			t.Fatal()
		}
		req.Header.Add(AUTHORIZATION, BEARER+ha)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fail()
		}
		res_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fail()
		}
		res_string := string(res_bytes)
		assert.Contains(t, res_string, "next")
		assert.Contains(t, res_string, "called")
		assert.Contains(t, res_string, "token")
		assert.Contains(t, res_string, "woefhweh")
	})
}

func CreateMocks(t *testing.T) (*common.MockAuthClient, *gin.Context, *common.MockErrorHandler) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	gin_context, _ := gin.CreateTestContext(httptest.NewRecorder())
	return common.NewMockAuthClient(mockCtrl), gin_context, common.NewMockErrorHandler(mockCtrl)
}
