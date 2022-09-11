package handler

import (
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
)

type ErrorHandlerImpl struct {
	WebContext common.WebContext
	Logger     common.Logger
}

type AppError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (eh ErrorHandlerImpl) HandleAppError(someError error, messege string, code int) {
	errObj := AppError{Error: someError.Error(), Message: messege}
	eh.Logger.Printf("%v\n", errObj)
	eh.WebContext.JSON(code, errObj)
}

func Create(todoRepository common.TodoRepository, errorHandler common.ErrorHandler) func(common.WebContext) {
	return func(ctx common.WebContext) {
		token, ok := ctx.Get(middleware.AuthToken)
		if !ok {
			errorHandler.HandleAppError(middleware.ErrNoUIDinToken, "", http.StatusUnauthorized)
		} else {
		var json model.Todo
		if err := ctx.ShouldBindJSON(&json); err != nil {
			errorHandler.HandleAppError(err,
				"", http.StatusBadRequest)
			return
		}
		token := token.(*auth.Token)
		err := todoRepository.Create(&json, token.UID)
		if err != nil {
			errorHandler.HandleAppError(err,
				"", http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, json)
		}
	}
	}
}