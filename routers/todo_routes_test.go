package routers

import (
	"reflect"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/controllers"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
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
	create := controllers.Create(todoRepositoryMock, errorHandlerMock)
	routerMock.EXPECT().POST("/todos", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(create).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getAll := controllers.GetAll(todoRepositoryMock, errorHandlerMock)
	routerMock.EXPECT().GET("/todos", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(getAll).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getById := controllers.GetById(todoRepositoryMock, errorHandlerMock, uuid.Parse)
	routerMock.EXPECT().GET("/todos/:id", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(getById).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getAllByUserId := controllers.GetAllByUserId(todoRepositoryMock, errorHandlerMock, uuid.Parse)
	routerMock.EXPECT().GET("/todos/users/:id", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(getAllByUserId).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	update := controllers.Update(todoRepositoryMock, errorHandlerMock)
	routerMock.EXPECT().PUT("/todos", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(update).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	delete := controllers.Delete(todoRepositoryMock, errorHandlerMock, uuid.Parse)
	routerMock.EXPECT().DELETE("/todos/:id", gomock.Any()).Do(func(path string, handler func(common.WebContext)) {
		assert.Equal(t, reflect.ValueOf(delete).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	SetTodoRoutes(routerMock, todoRepositoryMock, errorHandlerMock, firebaseAuthClientMock)
}
