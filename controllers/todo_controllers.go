package controllers

import (
	"errors"
	"net/http"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrParseIsNil error = errors.New("parse is nil")


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
