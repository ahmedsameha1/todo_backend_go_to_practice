// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ahmedsameha1/todo_backend_go_to_practice/common (interfaces: WebContext)

// Package common is a generated GoMock package.
package common

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWebContext is a mock of WebContext interface.
type MockWebContext struct {
	ctrl     *gomock.Controller
	recorder *MockWebContextMockRecorder
}

// MockWebContextMockRecorder is the mock recorder for MockWebContext.
type MockWebContextMockRecorder struct {
	mock *MockWebContext
}

// NewMockWebContext creates a new mock instance.
func NewMockWebContext(ctrl *gomock.Controller) *MockWebContext {
	mock := &MockWebContext{ctrl: ctrl}
	mock.recorder = &MockWebContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebContext) EXPECT() *MockWebContextMockRecorder {
	return m.recorder
}

// JSON mocks base method.
func (m *MockWebContext) JSON(arg0 int, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "JSON", arg0, arg1)
}

// JSON indicates an expected call of JSON.
func (mr *MockWebContextMockRecorder) JSON(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JSON", reflect.TypeOf((*MockWebContext)(nil).JSON), arg0, arg1)
}

// Param mocks base method.
func (m *MockWebContext) Param(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Param", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Param indicates an expected call of Param.
func (mr *MockWebContextMockRecorder) Param(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Param", reflect.TypeOf((*MockWebContext)(nil).Param), arg0)
}

// ShouldBindJSON mocks base method.
func (m *MockWebContext) ShouldBindJSON(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldBindJSON", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ShouldBindJSON indicates an expected call of ShouldBindJSON.
func (mr *MockWebContextMockRecorder) ShouldBindJSON(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldBindJSON", reflect.TypeOf((*MockWebContext)(nil).ShouldBindJSON), arg0)
}