package routers

import "github.com/gin-gonic/gin"

func SetTodoRoutes(poster router) router {
	poster.POST("/todos", func(ctx *gin.Context) {})
	poster.PUT("/todos/{id}", func(ctx *gin.Context) {})
	poster.GET("/todos", func(ctx *gin.Context) {})
	poster.GET("/todos/{id}", func(ctx *gin.Context) {})
	poster.GET("/todos/users/{id}", func(ctx *gin.Context) {})
	poster.DELETE("/todos/{id}", func(ctx *gin.Context) {})
	return poster
}
