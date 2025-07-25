// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repo is a generated GoMock package.
package repo

import (
	reflect "reflect"
	entities "salesapi/internal/sale/entities"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetTopProducts mocks base method.
func (m *MockRepository) GetTopProducts(req entities.GetTopProductsRequest) ([]entities.TopProduct, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopProducts", req)
	ret0, _ := ret[0].([]entities.TopProduct)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopProducts indicates an expected call of GetTopProducts.
func (mr *MockRepositoryMockRecorder) GetTopProducts(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopProducts", reflect.TypeOf((*MockRepository)(nil).GetTopProducts), req)
}
