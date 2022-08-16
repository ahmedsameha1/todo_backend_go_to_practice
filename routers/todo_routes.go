package routers

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
)

func SetTodoRoutes(router common.Router) common.Router {
	router.POST("/todos", func(ctx *gin.Context) {})	//
	router.GET("/todos", func(ctx *gin.Context) {})		//
	router.GET("/todos/{id}", func(ctx *gin.Context) {})
	router.GET("/todos/users/{id}", func(ctx *gin.Context) {})
	router.PUT("/todos/{id}", func(ctx *gin.Context) {})
	router.DELETE("/todos/{id}", func(ctx *gin.Context) {})
	return router
}
