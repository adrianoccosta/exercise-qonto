// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/adrianoccosta/exercise-qonto/internal/services/transfersvc (interfaces: TransferService)

// Package mockservice is a generated GoMock package.
package mockservice

import (
	reflect "reflect"

	domain "github.com/adrianoccosta/exercise-qonto/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockTransferService is a mock of TransferService interface.
type MockTransferService struct {
	ctrl     *gomock.Controller
	recorder *MockTransferServiceMockRecorder
}

// MockTransferServiceMockRecorder is the mock recorder for MockTransferService.
type MockTransferServiceMockRecorder struct {
	mock *MockTransferService
}

// NewMockTransferService creates a new mock instance.
func NewMockTransferService(ctrl *gomock.Controller) *MockTransferService {
	mock := &MockTransferService{ctrl: ctrl}
	mock.recorder = &MockTransferServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferService) EXPECT() *MockTransferServiceMockRecorder {
	return m.recorder
}

// BulkTransfer mocks base method.
func (m *MockTransferService) BulkTransfer(arg0 domain.BulkTransfer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkTransfer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkTransfer indicates an expected call of BulkTransfer.
func (mr *MockTransferServiceMockRecorder) BulkTransfer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkTransfer", reflect.TypeOf((*MockTransferService)(nil).BulkTransfer), arg0)
}