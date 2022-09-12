package routers

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/controllers"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/handler"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
	"github.com/google/uuid"
)

func SetTodoRoutes(router common.Router, todoRepository common.TodoRepository,
	errorHandler common.ErrorHandler, authClient common.AuthClient) common.Router {
	router.Use(middleware.GetAuthMiddleware(authClient, errorHandler))
	router.POST("/todos", handler.Create(todoRepository, errorHandler))
	router.GET("/todos", handler.GetAll(todoRepository, errorHandler))
	router.GET("/todos/:id", handler.GetById(todoRepository, errorHandler, uuid.Parse))
	router.PUT("/todos", handler.Update(todoRepository, errorHandler))
	router.DELETE("/todos/:id", controllers.Delete(todoRepository, errorHandler, uuid.Parse))
	return router
}
