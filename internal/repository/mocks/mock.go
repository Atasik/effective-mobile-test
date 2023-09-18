// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	domain "fio/internal/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPersonRepo is a mock of PersonRepo interface.
type MockPersonRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPersonRepoMockRecorder
}

// MockPersonRepoMockRecorder is the mock recorder for MockPersonRepo.
type MockPersonRepoMockRecorder struct {
	mock *MockPersonRepo
}

// NewMockPersonRepo creates a new mock instance.
func NewMockPersonRepo(ctrl *gomock.Controller) *MockPersonRepo {
	mock := &MockPersonRepo{ctrl: ctrl}
	mock.recorder = &MockPersonRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersonRepo) EXPECT() *MockPersonRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockPersonRepo) Add(person domain.Person) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", person)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockPersonRepoMockRecorder) Add(person interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockPersonRepo)(nil).Add), person)
}

// Delete mocks base method.
func (m *MockPersonRepo) Delete(personID int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", personID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockPersonRepoMockRecorder) Delete(personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPersonRepo)(nil).Delete), personID)
}

// GetAll mocks base method.
func (m *MockPersonRepo) GetAll(opts domain.PersonsQuery) ([]domain.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", opts)
	ret0, _ := ret[0].([]domain.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPersonRepoMockRecorder) GetAll(opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPersonRepo)(nil).GetAll), opts)
}

// Update mocks base method.
func (m *MockPersonRepo) Update(personID int, input domain.UpdatePersonInput) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", personID, input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockPersonRepoMockRecorder) Update(personID, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPersonRepo)(nil).Update), personID, input)
}