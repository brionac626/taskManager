// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/repository.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/repository.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	models "github.com/brionac626/taskManager/models"
	gomock "go.uber.org/mock/gomock"
)

// MockTaskManager is a mock of TaskManager interface.
type MockTaskManager struct {
	ctrl     *gomock.Controller
	recorder *MockTaskManagerMockRecorder
	isgomock struct{}
}

// MockTaskManagerMockRecorder is the mock recorder for MockTaskManager.
type MockTaskManagerMockRecorder struct {
	mock *MockTaskManager
}

// NewMockTaskManager creates a new mock instance.
func NewMockTaskManager(ctrl *gomock.Controller) *MockTaskManager {
	mock := &MockTaskManager{ctrl: ctrl}
	mock.recorder = &MockTaskManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskManager) EXPECT() *MockTaskManagerMockRecorder {
	return m.recorder
}

// CreateTasks mocks base method.
func (m *MockTaskManager) CreateTasks(ctx context.Context, tasks []models.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTasks", ctx, tasks)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTasks indicates an expected call of CreateTasks.
func (mr *MockTaskManagerMockRecorder) CreateTasks(ctx, tasks any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTasks", reflect.TypeOf((*MockTaskManager)(nil).CreateTasks), ctx, tasks)
}

// DeleteTask mocks base method.
func (m *MockTaskManager) DeleteTask(ctx context.Context, taskID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskManagerMockRecorder) DeleteTask(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskManager)(nil).DeleteTask), ctx, taskID)
}

// GetTasks mocks base method.
func (m *MockTaskManager) GetTasks(ctx context.Context) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasks", ctx)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasks indicates an expected call of GetTasks.
func (mr *MockTaskManagerMockRecorder) GetTasks(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasks", reflect.TypeOf((*MockTaskManager)(nil).GetTasks), ctx)
}

// UpdateTask mocks base method.
func (m *MockTaskManager) UpdateTask(ctx context.Context, taskID string, name *string, status *int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, taskID, name, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskManagerMockRecorder) UpdateTask(ctx, taskID, name, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskManager)(nil).UpdateTask), ctx, taskID, name, status)
}
