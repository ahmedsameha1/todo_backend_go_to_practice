package common

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestDisplayAppError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	loggerMock := NewMockLogger(mockCtrl)
	ginContextMock := NewMockWebContext(mockCtrl)
	handlerErr := errors.New("handlerErr1")
	message := "message1"
	code := 404
	errObj := AppError{Error: handlerErr.Error(), Message: message}
	ginContextMock.EXPECT().JSON(code, errObj)
	loggerMock.EXPECT().Printf("%v\n", errObj)
	SendBackAnAppError(ginContextMock, loggerMock, handlerErr, message, code)
}
