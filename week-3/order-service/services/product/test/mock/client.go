// Code generated by MockGen. DO NOT EDIT.
// Source: repository/aws/client.go
//
// Generated by this command:
//
//	mockgen -source repository/aws/client.go -destination test/mock/client.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAWSClient is a mock of AWSClient interface.
type MockAWSClient struct {
	ctrl     *gomock.Controller
	recorder *MockAWSClientMockRecorder
}

// MockAWSClientMockRecorder is the mock recorder for MockAWSClient.
type MockAWSClientMockRecorder struct {
	mock *MockAWSClient
}

// NewMockAWSClient creates a new mock instance.
func NewMockAWSClient(ctrl *gomock.Controller) *MockAWSClient {
	mock := &MockAWSClient{ctrl: ctrl}
	mock.recorder = &MockAWSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAWSClient) EXPECT() *MockAWSClientMockRecorder {
	return m.recorder
}

// DeleteImage mocks base method.
func (m *MockAWSClient) DeleteImage(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImage indicates an expected call of DeleteImage.
func (mr *MockAWSClientMockRecorder) DeleteImage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockAWSClient)(nil).DeleteImage), arg0, arg1)
}

// SaveImage mocks base method.
func (m *MockAWSClient) SaveImage(arg0 context.Context, arg1 *[]byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveImage", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveImage indicates an expected call of SaveImage.
func (mr *MockAWSClientMockRecorder) SaveImage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveImage", reflect.TypeOf((*MockAWSClient)(nil).SaveImage), arg0, arg1)
}