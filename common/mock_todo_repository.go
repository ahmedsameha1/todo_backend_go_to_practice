// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ahmedsameha1/todo_backend_go_to_practice/common (interfaces: TodoRepository)

// Package common is a generated GoMock package.
package common

import (
	reflect "reflect"

	model "github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockTodoRepository is a mock of TodoRepository interface.
type MockTodoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTodoRepositoryMockRecorder
}

// MockTodoRepositoryMockRecorder is the mock recorder for MockTodoRepository.
type MockTodoRepositoryMockRecorder struct {
	mock *MockTodoRepository
}

// NewMockTodoRepository creates a new mock instance.
func NewMockTodoRepository(ctrl *gomock.Controller) *MockTodoRepository {
	mock := &MockTodoRepository{ctrl: ctrl}
	mock.recorder = &MockTodoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoRepository) EXPECT() *MockTodoRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTodoRepository) Create(arg0 *model.Todo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTodoRepositoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTodoRepository)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockTodoRepository) Delete(arg0 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTodoRepositoryMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTodoRepository)(nil).Delete), arg0)
}

// GetAll mocks base method.
func (m *MockTodoRepository) GetAll() ([]model.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTodoRepositoryMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTodoRepository)(nil).GetAll))
}

// GetAllByUserId mocks base method.
func (m *MockTodoRepository) GetAllByUserId(arg0 uuid.UUID) ([]model.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserId", arg0)
	ret0, _ := ret[0].([]model.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUserId indicates an expected call of GetAllByUserId.
func (mr *MockTodoRepositoryMockRecorder) GetAllByUserId(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserId", reflect.TypeOf((*MockTodoRepository)(nil).GetAllByUserId), arg0)
}

// GetById mocks base method.
func (m *MockTodoRepository) GetById(arg0 uuid.UUID) (*model.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", arg0)
	ret0, _ := ret[0].(*model.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockTodoRepositoryMockRecorder) GetById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockTodoRepository)(nil).GetById), arg0)
}

// Update mocks base method.
func (m *MockTodoRepository) Update(arg0 *model.Todo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTodoRepositoryMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTodoRepository)(nil).Update), arg0)
}
