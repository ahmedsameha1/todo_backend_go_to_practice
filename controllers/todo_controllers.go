package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrParseIsNil error = errors.New("parse is nil")
var logger *log.Logger = log.Default()

func Create(todoRepository common.TodoRepository, errorHandler common.ErrorHandler) func(common.WebContext) {
	return func(ctx common.WebContext) {
		var json model.Todo
		if err := ctx.ShouldBindJSON(&json); err != nil {
			errorHandler.HandleAppError(err,
				"", http.StatusBadRequest)
			return
		}
		err := todoRepository.Create(&json)
		if err != nil {
			errorHandler.HandleAppError(err,
				"", http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, json)
		}
	}
}

func Update(todoRepository common.TodoRepository, errorHandler common.ErrorHandler,
	parse func(string) (uuid.UUID, error)) func(common.WebContext) {
	return func(ctx common.WebContext) {
		if parse != nil {
			_, err := parse(ctx.Param("id"))
			if err != nil {
				errorHandler.HandleAppError(err, "", http.StatusBadRequest)
			} else {
				var todo model.Todo
				err := ctx.ShouldBindJSON(&todo)
				if err != nil {
					errorHandler.HandleAppError(err, "", http.StatusBadRequest)
				} else {
					err := todoRepository.Update(&todo)
					if err != nil {
						errorHandler.HandleAppError(err, "", http.StatusInternalServerError)
					} else {
						ctx.JSON(http.StatusNoContent, gin.H{})
					}
				}
			}
		} else {
			errorHandler.HandleAppError(ErrParseIsNil, "", http.StatusInternalServerError)
		}
	}
}

func GetAll(todoRepository common.TodoRepository, errorHandler common.ErrorHandler) func(common.WebContext) {
	return func(ctx common.WebContext) {
		if todos, err := todoRepository.GetAll(); err != nil {
			errorHandler.HandleAppError(err, "", http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, todos)
		}
	}
}

func GetById(todoRepository common.TodoRepository, errorHandler common.ErrorHandler,
	parse func(string) (uuid.UUID, error)) func(common.WebContext) {
	return func(ctx common.WebContext) {
		if parse != nil {
			id, err := parse(ctx.Param("id"))
			if err != nil {
				errorHandler.HandleAppError(err, "", http.StatusBadRequest)
			} else {
				todo, err := todoRepository.GetById(id)
				if err != nil {
					errorHandler.HandleAppError(err, "", http.StatusInternalServerError)
				} else {
					ctx.JSON(http.StatusOK, todo)
				}
			}
		} else {
			errorHandler.HandleAppError(ErrParseIsNil, "", http.StatusInternalServerError)
		}
	}
}

func GetAllByUserId(todoRepository common.TodoRepository, errorHandler common.ErrorHandler,
	parse func(string) (uuid.UUID, error)) func(common.WebContext) {
	return func(ctx common.WebContext) {
		if parse != nil {
			id, err := parse(ctx.Param("id"))
			if err != nil {
				errorHandler.HandleAppError(err, "", http.StatusBadRequest)
			} else {
				todos, err := todoRepository.GetAllByUserId(id)
				if err != nil {
					errorHandler.HandleAppError(err, "", http.StatusInternalServerError)
				} else {
					ctx.JSON(http.StatusOK, todos)
				}
			}
		} else {
			errorHandler.HandleAppError(ErrParseIsNil, "", http.StatusInternalServerError)
		}
	}
}

func Delete(todoRepository common.TodoRepository, errorHandler common.ErrorHandler,
	parse func(string) (uuid.UUID, error)) func(common.WebContext) {
	return func(ctx common.WebContext) {
		if parse != nil {
			id, err := parse(ctx.Param("id"))
			if err != nil {
				errorHandler.HandleAppError(err, "", http.StatusBadRequest)
			} else {
				err := todoRepository.Delete(id)
				if err != nil {
					errorHandler.HandleAppError(err, "", http.StatusInternalServerError)
				} else {
					ctx.JSON(http.StatusNoContent, gin.H{})
				}
			}
		} else {
			errorHandler.HandleAppError(ErrParseIsNil, "", http.StatusInternalServerError)
		}
	}
}
