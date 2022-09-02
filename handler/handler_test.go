package handler

import (
	"errors"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestHandleError(t *testing.T) {
	t.Run("Normal case", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		loggerMock := common.NewMockLogger(mockCtrl)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		handlerErr := errors.New("handlerErr1")
		message := "message1"
		code := 404
		errObj := common.AppError{Error: handlerErr.Error(), Message: message}
		loggerMock.EXPECT().Printf("%v\n", errObj)
		ginContextMock.EXPECT().JSON(code, errObj)
		errorHandlerImpl := ErrorHandlerImpl{WebContext: ginContextMock, Logger: loggerMock}
		errorHandlerImpl.HandleAppError(handlerErr, message, code)
	})

	t.Run("When WebContext or Logger is nil, I trust that the app will panic!!", func(t *testing.T) {})
}
