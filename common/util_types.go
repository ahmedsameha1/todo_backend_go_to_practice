package common

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

type Router interface {
	POST(relativePath string, handlers ...func(WebContext)) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

type Logger interface {
	Printf(format string, v ...any)
}

type TodoRepository interface {
	Create(todo *model.Todo) error
	GetAll() ([]model.Todo, error)
	GetById(id uuid.UUID) (*model.Todo, error)
	GetAllByUserId(id uuid.UUID) ([]model.Todo, error)
	Update(todo *model.Todo) error
	Delete(id uuid.UUID) error
}

type DBPool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (DBRows, error)
}

type DBRows interface {
	Err() error
	Next() bool
	Scan(dest ...interface{}) (err error)
}

type ErrorHandler interface {
	HandleAppError(error, string, int)
}

type AuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type WebContext interface {
	JSON(code int, obj any)
	ShouldBindJSON(obj any) error
}