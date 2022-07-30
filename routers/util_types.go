package routers

import (
	"github.com/gin-gonic/gin"
)

type router interface {
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

type PosterMock struct {
	called int8
}

func (p *PosterMock) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.called++
	return nil
}

func (p *PosterMock) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.called++
	return nil
}

func (p *PosterMock) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.called++
	return nil
}

func (p *PosterMock) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	p.called++
	return nil
}

//const handlerPlaceholder func(ctx *gin.Context) = func(ctx *gin.Context) {}
