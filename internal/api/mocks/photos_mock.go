// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/api/photos.go

// Package mock_api is a generated GoMock package.
package mock_api

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	photos "github.com/twk/skeleton-go-api/internal/photos"
)

// MockphotoService is a mock of photoService interface.
type MockphotoService struct {
	ctrl     *gomock.Controller
	recorder *MockphotoServiceMockRecorder
}

// MockphotoServiceMockRecorder is the mock recorder for MockphotoService.
type MockphotoServiceMockRecorder struct {
	mock *MockphotoService
}

// NewMockphotoService creates a new mock instance.
func NewMockphotoService(ctrl *gomock.Controller) *MockphotoService {
	mock := &MockphotoService{ctrl: ctrl}
	mock.recorder = &MockphotoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockphotoService) EXPECT() *MockphotoServiceMockRecorder {
	return m.recorder
}

// GetPhotos mocks base method.
func (m *MockphotoService) GetPhotos(ctx context.Context, albumID int) (*photos.Photo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPhotos", ctx, albumID)
	ret0, _ := ret[0].(*photos.Photo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPhotos indicates an expected call of GetPhotos.
func (mr *MockphotoServiceMockRecorder) GetPhotos(ctx, albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPhotos", reflect.TypeOf((*MockphotoService)(nil).GetPhotos), ctx, albumID)
}
