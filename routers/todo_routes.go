package routers

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/controllers"
)

func SetTodoRoutes(router common.Router, todoRepository common.TodoRepository, errorHandler common.ErrorHandler) common.Router {
	router.POST("/todos", controllers.Create(todoRepository, errorHandler))
	router.GET("/todos", controllers.GetAll(todoRepository, errorHandler))
	router.GET("/todos/{id}", controllers.GetById(todoRepository))
	router.GET("/todos/users/{id}", controllers.GetAllByUserId(todoRepository))
	router.PUT("/todos/{id}", controllers.Update(todoRepository))
	router.DELETE("/todos/{id}", controllers.Delete(todoRepository))
	return router
}
