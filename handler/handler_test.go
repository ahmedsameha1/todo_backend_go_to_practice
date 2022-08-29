package handler

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	t.Run("Normal case", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		mockCtrl := gomock.NewController(t)
		loggerMock := common.NewMockLogger(mockCtrl)
		handlerErr := errors.New("handlerErr1")
		message := "message1"
		code := 404
		errObj := common.AppError{Error: handlerErr.Error(), Message: message}
		loggerMock.EXPECT().Printf("%v\n", errObj)
		errorHandlerImpl := ErrorHandlerImpl{GinContext: ctx, Logger: loggerMock}
		errorHandlerImpl.HandleAppError(handlerErr, message, code)
		if w.Code != int(code) {
			t.Errorf("Expects %d, but got %d", int(code), w.Code)
		}
		var got gin.H
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		want := gin.H{"error": handlerErr.Error(), "message": message}
		assert.Equal(t, want, got)
	})

	t.Run("GinContext is nil", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		loggerMock := common.NewMockLogger(mockCtrl)
		handlerErr := errors.New("handlerErr1")
		message := "message1"
		code := 404
		errObj := common.AppError{Error: handlerErr.Error(), Message: message}
		loggerMock.EXPECT().Printf("%v\n", errObj)
		errorHandlerImpl := ErrorHandlerImpl{GinContext: nil, Logger: loggerMock}
		assert.Panics(t, func() { errorHandlerImpl.HandleAppError(handlerErr, message, code) })
	})

	t.Run("GinContext is nil", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		handlerErr := errors.New("handlerErr1")
		message := "message1"
		code := 404
		errorHandlerImpl := ErrorHandlerImpl{GinContext: ctx, Logger: nil}
		assert.Panics(t, func() { errorHandlerImpl.HandleAppError(handlerErr, message, code) })
	})
}
