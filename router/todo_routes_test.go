package router

import (
	"reflect"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/handler"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSetTodoRoutes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	routerMock := common.NewMockRouter(mockCtrl)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
	firebaseAuthClientMock := common.NewMockAuthClient(mockCtrl)
	authMiddleware := middleware.GetAuthMiddleware(firebaseAuthClientMock, errorHandlerMock)
	routerMock.EXPECT().Use(gomock.Any()).Do(func(handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(authMiddleware).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	create := handler.Create(todoRepositoryMock, errorHandlerMock)
	routerMock.EXPECT().POST("/todos", gomock.Any()).Do(func(path string, handler func(*gin.Context)) {
		assert.Equal(t, reflect.ValueOf(create).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getAll := handler.GetAll(todoRepositoryMock, errorHandlerMock)
	routerMock.EXPECT().GET("/todos", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(getAll).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getById := handler.GetById(todoRepositoryMock, errorHandlerMock, uuid.Parse)
	routerMock.EXPECT().GET("/todos/:id", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(getById).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	update := handler.Update(todoRepositoryMock, errorHandlerMock)
	routerMock.EXPECT().PUT("/todos", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(update).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	delete := handler.Delete(todoRepositoryMock, errorHandlerMock, uuid.Parse)
	routerMock.EXPECT().DELETE("/todos/:id", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(delete).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	SetTodoRoutes(routerMock, todoRepositoryMock, errorHandlerMock, firebaseAuthClientMock)
}
