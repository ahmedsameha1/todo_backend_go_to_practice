package common

import (
	"context"
	"errors"

	"firebase.google.com/go/v4/auth"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
)

var ErrError error = errors.New("an error")

type Router interface {
	Use(middleware ...func(WebContext)) gin.IRoutes
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

type Logger interface {
	Printf(format string, v ...any)
}

type TodoRepository interface {
	Create(todo *model.Todo, userId string) error
	GetAll(userId string) ([]model.Todo, error)
	GetById(id string, userId string) (*model.Todo, error)
	Update(todo *model.Todo, userId string) error
	Delete(id string, userId string) error
}

type ErrorHandler interface {
	HandleAppError(WebContext, error, int)
}

type AuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type WebContext interface {
	JSON(code int, obj any)
	ShouldBindJSON(obj any) error
	Param(key string) string
	GetHeader(key string) string
	Set(key string, value any)
	Get(key string) (value any, exists bool)
	Next()
}
