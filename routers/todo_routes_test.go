package routers

import (
	"reflect"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/controllers"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetTodoRoutes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	routerMock := common.NewMockRouter(mockCtrl)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	create := controllers.Create(todoRepositoryMock)
	routerMock.EXPECT().POST("/todos", gomock.Any()).Do(func(path string, handler gin.HandlerFunc) {
		assert.Equal(t, reflect.ValueOf(create).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getAll := controllers.GetAll(todoRepositoryMock)
	routerMock.EXPECT().GET("/todos", gomock.Any()).Do(func(path string, handler gin.HandlerFunc) {
		assert.Equal(t, reflect.ValueOf(getAll).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getById := controllers.GetById(todoRepositoryMock)
	routerMock.EXPECT().GET("/todos/{id}", gomock.Any()).Do(func(path string, handler gin.HandlerFunc) {
		assert.Equal(t, reflect.ValueOf(getById).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	getAllByUserId := controllers.GetAllByUserId(todoRepositoryMock)
	routerMock.EXPECT().GET("/todos/users/{id}", gomock.Any()).Do(func(path string, handler gin.HandlerFunc) {
		assert.Equal(t, reflect.ValueOf(getAllByUserId).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	update := controllers.Update(todoRepositoryMock)
	routerMock.EXPECT().PUT("/todos/{id}", gomock.Any()).Do(func(path string, handler gin.HandlerFunc) {
		assert.Equal(t, reflect.ValueOf(update).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	delete := controllers.Delete(todoRepositoryMock)
	routerMock.EXPECT().DELETE("/todos/{id}", gomock.Any()).Do(func(path string, handler gin.HandlerFunc) {
		assert.Equal(t, reflect.ValueOf(delete).Pointer(), reflect.ValueOf(handler).Pointer())
	})
	SetTodoRoutes(routerMock, todoRepositoryMock)
}
