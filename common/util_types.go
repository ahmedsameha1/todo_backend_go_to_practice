package common

import (
	"github.com/gin-gonic/gin"
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