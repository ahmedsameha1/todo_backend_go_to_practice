package common

import (
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Router interface {
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

type PosterMock struct {
	Called int8
}

func (p *PosterMock) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.Called++
	return nil
}

func (p *PosterMock) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.Called++
	return nil
}

func (p *PosterMock) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.Called++
	return nil
}

func (p *PosterMock) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.Called++
	return nil
}

type Logger interface{
	Printf(format string, v ...any)
}

type LoggerMock struct {
	Called int
}

func (l *LoggerMock) Printf(format string, v ...any) {
	l.Called++
}

type TodoRepository interface {
	Create(*model.Todo) (*model.Todo, error)
	GetAll() ([]model.Todo, error)
	GetById(id uuid.UUID) (*model.Todo, error)
	GetAllByUserId(id uuid.UUID) ([]model.Todo, error)
}
