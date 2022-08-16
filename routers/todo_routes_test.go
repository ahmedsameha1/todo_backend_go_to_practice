package routers

import (
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/golang/mock/gomock"
)

func TestSetTodoRoutes(t *testing.T){
	mockCtrl := gomock.NewController(t)
	routerMock := common.NewMockRouter(mockCtrl)
	routerMock.EXPECT().POST("/todos", gomock.Any())
	routerMock.EXPECT().GET("/todos", gomock.Any())
	routerMock.EXPECT().GET("/todos/{id}", gomock.Any())
	routerMock.EXPECT().GET("/todos/users/{id}", gomock.Any())
	routerMock.EXPECT().PUT("/todos/{id}", gomock.Any()) 
	routerMock.EXPECT().DELETE("/todos/{id}", gomock.Any()) 
	SetTodoRoutes(routerMock)
}