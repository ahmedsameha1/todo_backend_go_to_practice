// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ahmedsameha1/todo_backend_go_to_practice/common (interfaces: Router)

// Package common is a generated GoMock package.
package common

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockRouter is a mock of Router interface.
type MockRouter struct {
	ctrl     *gomock.Controller
	recorder *MockRouterMockRecorder
}

// MockRouterMockRecorder is the mock recorder for MockRouter.
type MockRouterMockRecorder struct {
	mock *MockRouter
}

// NewMockRouter creates a new mock instance.
func NewMockRouter(ctrl *gomock.Controller) *MockRouter {
	mock := &MockRouter{ctrl: ctrl}
	mock.recorder = &MockRouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRouter) EXPECT() *MockRouterMockRecorder {
	return m.recorder
}

// DELETE mocks base method.
func (m *MockRouter) DELETE(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DELETE", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// DELETE indicates an expected call of DELETE.
func (mr *MockRouterMockRecorder) DELETE(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DELETE", reflect.TypeOf((*MockRouter)(nil).DELETE), varargs...)
}

// GET mocks base method.
func (m *MockRouter) GET(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GET", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// GET indicates an expected call of GET.
func (mr *MockRouterMockRecorder) GET(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GET", reflect.TypeOf((*MockRouter)(nil).GET), varargs...)
}

// POST mocks base method.
func (m *MockRouter) POST(arg0 string, arg1 ...func(WebContext)) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "POST", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// POST indicates an expected call of POST.
func (mr *MockRouterMockRecorder) POST(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POST", reflect.TypeOf((*MockRouter)(nil).POST), varargs...)
}

// PUT mocks base method.
func (m *MockRouter) PUT(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PUT", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// PUT indicates an expected call of PUT.
func (mr *MockRouterMockRecorder) PUT(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PUT", reflect.TypeOf((*MockRouter)(nil).PUT), varargs...)
}
