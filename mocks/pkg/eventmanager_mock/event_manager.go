// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/liatrio/rode/pkg/eventmanager (interfaces: EventManager)

// Package eventmanager_mock is a generated GoMock package.
package eventmanager_mock

import (
	gomock "github.com/golang/mock/gomock"
	grafeas_go_proto "github.com/grafeas/grafeas/proto/v1beta1/grafeas_go_proto"
	reflect "reflect"
)

// MockEventManager is a mock of EventManager interface
type MockEventManager struct {
	ctrl     *gomock.Controller
	recorder *MockEventManagerMockRecorder
}

// MockEventManagerMockRecorder is the mock recorder for MockEventManager
type MockEventManagerMockRecorder struct {
	mock *MockEventManager
}

// NewMockEventManager creates a new mock instance
func NewMockEventManager(ctrl *gomock.Controller) *MockEventManager {
	mock := &MockEventManager{ctrl: ctrl}
	mock.recorder = &MockEventManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEventManager) EXPECT() *MockEventManagerMockRecorder {
	return m.recorder
}

// Initialize mocks base method
func (m *MockEventManager) Initialize(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize
func (mr *MockEventManagerMockRecorder) Initialize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockEventManager)(nil).Initialize), arg0)
}

// Publish mocks base method
func (m *MockEventManager) Publish(arg0 string, arg1 *grafeas_go_proto.Occurrence) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockEventManagerMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockEventManager)(nil).Publish), arg0, arg1)
}

// Subscribe mocks base method
func (m *MockEventManager) Subscribe(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockEventManagerMockRecorder) Subscribe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockEventManager)(nil).Subscribe), arg0)
}

// Unsubscribe mocks base method
func (m *MockEventManager) Unsubscribe(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe
func (mr *MockEventManagerMockRecorder) Unsubscribe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockEventManager)(nil).Unsubscribe), arg0)
}