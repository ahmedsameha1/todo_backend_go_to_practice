package handler

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
)

type ErrorHandlerImpl struct {
	GinContext *gin.Context
	Logger     common.Logger
}

func (eh ErrorHandlerImpl) HandleAppError(someError error, messege string, code int) {
	errObj := common.AppError{Error: someError.Error(), Message: messege}
	eh.Logger.Printf("%v\n", errObj)
	eh.GinContext.JSON(int(code), errObj)
}
