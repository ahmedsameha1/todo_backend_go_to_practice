package handler

import (
	"errors"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrParseIsNil error = errors.New("parse is nil")

type ErrorHandlerImpl struct {
	Logger     common.Logger
}

type AppError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func (eh ErrorHandlerImpl) HandleAppError(webContext common.WebContext, someError error, messege string, code int) {
	errObj := AppError{Error: someError.Error(), Message: messege}
	eh.Logger.Printf("%v\n", errObj)
	webContext.JSON(code, errObj)
}

func Create(todoRepository common.TodoRepository, errorHandler common.ErrorHandler) func(common.WebContext) {
	return func(ctx common.WebContext) {
		token, ok := ctx.Get(middleware.AuthToken)
		if !ok {
			errorHandler.HandleAppError(ctx, middleware.ErrNoUID, http.StatusUnauthorized)
		} else {
		var json model.Todo
		if err := ctx.ShouldBindJSON(&json); err != nil {
			errorHandler.HandleAppError(ctx, err,
				http.StatusBadRequest)
			return
		}
		token := token.(*auth.Token)
		err := todoRepository.Create(&json, token.UID)
		if err != nil {
			errorHandler.HandleAppError(ctx, err,
				http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, json)
		}
	}
	}
}

func GetAll(todoRepository common.TodoRepository, errorHandler common.ErrorHandler) func(common.WebContext) {
	return func(ctx common.WebContext) {
		tokeN, ok := ctx.Get(middleware.AuthToken)
		if !ok {
			errorHandler.HandleAppError(ctx, middleware.ErrNoUID, http.StatusUnauthorized)
		} else {
			token := tokeN.(*auth.Token)
			if todos, err := todoRepository.GetAll(token.UID); err != nil {
				errorHandler.HandleAppError(ctx, err, http.StatusInternalServerError)
			} else {
				ctx.JSON(http.StatusOK, todos)
			}
		}
	}
}

func GetById(todoRepository common.TodoRepository, errorHandler common.ErrorHandler,
	parse func(string) (uuid.UUID, error)) func(common.WebContext) {
	return func(ctx common.WebContext) {
		token, ok := ctx.Get(middleware.AuthToken)
		if !ok {
			errorHandler.HandleAppError(ctx, middleware.ErrNoUID, http.StatusUnauthorized)
		} else {
			if parse != nil {
				id := ctx.Param("id")
				_, err := parse(id)
				if err != nil {
					errorHandler.HandleAppError(ctx, err, http.StatusBadRequest)
				} else {
					token := token.(*auth.Token)
					todo, err := todoRepository.GetById(id, token.UID)
					if err != nil {
						errorHandler.HandleAppError(ctx, err, http.StatusInternalServerError)
					} else {
						ctx.JSON(http.StatusOK, todo)
					}
				}
			} else {
				errorHandler.HandleAppError(ctx, ErrParseIsNil, http.StatusInternalServerError)
			}
		}
	}
}

func Update(todoRepository common.TodoRepository, errorHandler common.ErrorHandler) func(common.WebContext) {
	return func(ctx common.WebContext) {
		token, ok := ctx.Get(middleware.AuthToken)
		if !ok {
			errorHandler.HandleAppError(ctx, middleware.ErrNoUID, http.StatusUnauthorized)
		} else {
			var todo model.Todo
			err := ctx.ShouldBindJSON(&todo)
			if err != nil {
				errorHandler.HandleAppError(ctx, err, http.StatusBadRequest)
			} else {
				token := token.(*auth.Token)
				err := todoRepository.Update(&todo, token.UID)
				if err != nil {
					errorHandler.HandleAppError(ctx, err, http.StatusInternalServerError)
				} else {
					ctx.JSON(http.StatusNoContent, gin.H{})
				}
			}
		}
	}
}

func Delete(todoRepository common.TodoRepository, errorHandler common.ErrorHandler,
	parse func(string) (uuid.UUID, error)) func(common.WebContext) {
	return func(ctx common.WebContext) {
		token, ok := ctx.Get(middleware.AuthToken)
		if !ok {
			errorHandler.HandleAppError(ctx, middleware.ErrNoUID, http.StatusUnauthorized)
		} else {
			if parse != nil {
				id := ctx.Param("id")
				_, err := parse(id)
				if err != nil {
					errorHandler.HandleAppError(ctx, err, http.StatusBadRequest)
				} else {
					token := token.(*auth.Token)
					err := todoRepository.Delete(id, token.UID)
					if err != nil {
						errorHandler.HandleAppError(ctx, err, http.StatusInternalServerError)
					} else {
						ctx.JSON(http.StatusNoContent, gin.H{})
					}
				}
			} else {
				errorHandler.HandleAppError(ctx, ErrParseIsNil, http.StatusInternalServerError)
			}
		}
	}
}