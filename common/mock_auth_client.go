// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ahmedsameha1/todo_backend_go_to_practice/common (interfaces: AuthClient)

// Package common is a generated GoMock package.
package common

import (
	context "context"
	reflect "reflect"

	auth "firebase.google.com/go/v4/auth"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthClient is a mock of AuthClient interface.
type MockAuthClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthClientMockRecorder
}

// MockAuthClientMockRecorder is the mock recorder for MockAuthClient.
type MockAuthClientMockRecorder struct {
	mock *MockAuthClient
}

// NewMockAuthClient creates a new mock instance.
func NewMockAuthClient(ctrl *gomock.Controller) *MockAuthClient {
	mock := &MockAuthClient{ctrl: ctrl}
	mock.recorder = &MockAuthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthClient) EXPECT() *MockAuthClientMockRecorder {
	return m.recorder
}

// VerifyIDToken mocks base method.
func (m *MockAuthClient) VerifyIDToken(arg0 context.Context, arg1 string) (*auth.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyIDToken", arg0, arg1)
	ret0, _ := ret[0].(*auth.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyIDToken indicates an expected call of VerifyIDToken.
func (mr *MockAuthClientMockRecorder) VerifyIDToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyIDToken", reflect.TypeOf((*MockAuthClient)(nil).VerifyIDToken), arg0, arg1)
}
