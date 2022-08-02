package routers

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/gin-gonic/gin"
)

func SetTodoRoutes(poster common.Router) common.Router {
	poster.POST("/todos", func(ctx *gin.Context) {})
	poster.PUT("/todos/{id}", func(ctx *gin.Context) {})
	poster.GET("/todos", func(ctx *gin.Context) {})
	poster.GET("/todos/{id}", func(ctx *gin.Context) {})
	poster.GET("/todos/users/{id}", func(ctx *gin.Context) {})
	poster.DELETE("/todos/{id}", func(ctx *gin.Context) {})
	return poster
}
