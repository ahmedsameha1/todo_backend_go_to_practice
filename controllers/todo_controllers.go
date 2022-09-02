package controllers

import (
	"log"
	"net/http"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func Update(todoRepository common.TodoRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			common.SendBackAnAppError(ctx, logger, err, "", http.StatusBadRequest)
		} else {
			var todo model.Todo
			err := ctx.ShouldBindJSON(&todo)
			if err != nil {
				common.SendBackAnAppError(ctx, logger, err, "", http.StatusBadRequest)
			} else {
				err := todoRepository.Update(&todo)
				if err != nil {
					common.SendBackAnAppError(ctx, logger, err, "", http.StatusInternalServerError)
				} else {
					ctx.JSON(http.StatusNoContent, gin.H{})
				}
			}
		}
	}
}

func GetAll(todoRepository common.TodoRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if todos, err := todoRepository.GetAll(); err != nil {
			common.SendBackAnAppError(ctx, logger, err, "", http.StatusInternalServerError)
		} else {
			ctx.JSON(http.StatusOK, todos)
		}
	}
}

func GetById(todoRepository common.TodoRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			common.SendBackAnAppError(ctx, logger, err, "", http.StatusBadRequest)
		} else {
			todo, err := todoRepository.GetById(id)
			if err != nil {
				common.SendBackAnAppError(ctx, logger, err, "", http.StatusInternalServerError)
			} else {
				ctx.JSON(http.StatusOK, todo)
			}
		}
	}
}

func GetAllByUserId(todoRepository common.TodoRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			common.SendBackAnAppError(ctx, logger, err, "", http.StatusBadRequest)
		} else {
			todos, err := todoRepository.GetAllByUserId(id)
			if err != nil {
				common.SendBackAnAppError(ctx, logger, err, "", http.StatusInternalServerError)
			} else {
				ctx.JSON(http.StatusOK, todos)
			}
		}
	}
}

func Delete(todoRepository common.TodoRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			common.SendBackAnAppError(ctx, logger, err, "", http.StatusBadRequest)
		} else {
			err := todoRepository.Delete(id)
			if err != nil {
				common.SendBackAnAppError(ctx, logger, err, "", http.StatusInternalServerError)
			} else {
				ctx.JSON(http.StatusNoContent, gin.H{})
			}
		}
	}
}
