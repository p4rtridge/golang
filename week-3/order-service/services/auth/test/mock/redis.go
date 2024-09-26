// Code generated by MockGen. DO NOT EDIT.
// Source: repository/redis/store.go
//
// Generated by this command:
//
//	mockgen -source repository/redis/store.go -destination test/mock/redis.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTokenRepository is a mock of TokenRepository interface.
type MockTokenRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTokenRepositoryMockRecorder
}

// MockTokenRepositoryMockRecorder is the mock recorder for MockTokenRepository.
type MockTokenRepositoryMockRecorder struct {
	mock *MockTokenRepository
}

// NewMockTokenRepository creates a new mock instance.
func NewMockTokenRepository(ctrl *gomock.Controller) *MockTokenRepository {
	mock := &MockTokenRepository{ctrl: ctrl}
	mock.recorder = &MockTokenRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenRepository) EXPECT() *MockTokenRepositoryMockRecorder {
	return m.recorder
}

// DeleteAllRefreshToken mocks base method.
func (m *MockTokenRepository) DeleteAllRefreshToken(ctx context.Context, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllRefreshToken", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllRefreshToken indicates an expected call of DeleteAllRefreshToken.
func (mr *MockTokenRepositoryMockRecorder) DeleteAllRefreshToken(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllRefreshToken", reflect.TypeOf((*MockTokenRepository)(nil).DeleteAllRefreshToken), ctx, userID)
}

// DeleteRefreshToken mocks base method.
func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, userID int, deviceID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRefreshToken", ctx, userID, deviceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRefreshToken indicates an expected call of DeleteRefreshToken.
func (mr *MockTokenRepositoryMockRecorder) DeleteRefreshToken(ctx, userID, deviceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRefreshToken", reflect.TypeOf((*MockTokenRepository)(nil).DeleteRefreshToken), ctx, userID, deviceID)
}

// GetRefreshToken mocks base method.
func (m *MockTokenRepository) GetRefreshToken(ctx context.Context, userID int, deviceID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshToken", ctx, userID, deviceID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefreshToken indicates an expected call of GetRefreshToken.
func (mr *MockTokenRepositoryMockRecorder) GetRefreshToken(ctx, userID, deviceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshToken", reflect.TypeOf((*MockTokenRepository)(nil).GetRefreshToken), ctx, userID, deviceID)
}

// SetRefreshToken mocks base method.
func (m *MockTokenRepository) SetRefreshToken(ctx context.Context, userID int, deviceID, token string, expiration int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRefreshToken", ctx, userID, deviceID, token, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRefreshToken indicates an expected call of SetRefreshToken.
func (mr *MockTokenRepositoryMockRecorder) SetRefreshToken(ctx, userID, deviceID, token, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRefreshToken", reflect.TypeOf((*MockTokenRepository)(nil).SetRefreshToken), ctx, userID, deviceID, token, expiration)
}
