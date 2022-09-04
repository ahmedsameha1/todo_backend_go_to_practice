package routers

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/controllers"
	"github.com/google/uuid"
)

func SetTodoRoutes(router common.Router, todoRepository common.TodoRepository,
	errorHandler common.ErrorHandler) common.Router {
	router.POST("/todos", controllers.Create(todoRepository, errorHandler))
	router.GET("/todos", controllers.GetAll(todoRepository, errorHandler))
	router.GET("/todos/:id", controllers.GetById(todoRepository, errorHandler, uuid.Parse))
	router.GET("/todos/users/:id", controllers.GetAllByUserId(todoRepository, errorHandler, uuid.Parse))
	router.PUT("/todos/:id", controllers.Update(todoRepository, errorHandler, uuid.Parse))
	router.DELETE("/todos/:id", controllers.Delete(todoRepository))
	return router
}
