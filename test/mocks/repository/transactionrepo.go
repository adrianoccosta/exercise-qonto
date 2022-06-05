// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/adrianoccosta/exercise-qonto/internal/repository/transactionrepo (interfaces: TransactionRepository)

// Package mockrepository is a generated GoMock package.
package mockrepository

import (
	reflect "reflect"

	domain "github.com/adrianoccosta/exercise-qonto/internal/domain"
	transactionrepo "github.com/adrianoccosta/exercise-qonto/internal/repository/transactionrepo"
	gomock "github.com/golang/mock/gomock"
)

// MockTransactionRepository is a mock of TransactionRepository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTransactionRepository) Create(arg0 transactionrepo.Transaction) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTransactionRepositoryMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTransactionRepository)(nil).Create), arg0)
}

// Read mocks base method.
func (m *MockTransactionRepository) Read(arg0 uint) (transactionrepo.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(transactionrepo.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockTransactionRepositoryMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockTransactionRepository)(nil).Read), arg0)
}

// ReadByFilter mocks base method.
func (m *MockTransactionRepository) ReadByFilter(arg0 map[string]string) (domain.TransactionList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadByFilter", arg0)
	ret0, _ := ret[0].(domain.TransactionList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByFilter indicates an expected call of ReadByFilter.
func (mr *MockTransactionRepositoryMockRecorder) ReadByFilter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByFilter", reflect.TypeOf((*MockTransactionRepository)(nil).ReadByFilter), arg0)
}
