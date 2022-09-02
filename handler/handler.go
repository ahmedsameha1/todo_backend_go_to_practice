package handler

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
)

type ErrorHandlerImpl struct {
	WebContext common.WebContext
	Logger     common.Logger
}

func (eh ErrorHandlerImpl) HandleAppError(someError error, messege string, code int) {
	errObj := common.AppError{Error: someError.Error(), Message: messege}
	eh.Logger.Printf("%v\n", errObj)
	eh.WebContext.JSON(code, errObj)
}
