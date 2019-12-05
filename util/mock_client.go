// Code generated by MockGen. DO NOT EDIT.
// Source: util/client.go

// Package util is a generated GoMock package.
package util

import (
	gomock "github.com/golang/mock/gomock"
	client "github.com/longhorn/longhorn-manager/client"
	reflect "reflect"
)

// MockManagerClientInterface is a mock of ManagerClientInterface interface
type MockManagerClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockManagerClientInterfaceMockRecorder
}

// MockManagerClientInterfaceMockRecorder is the mock recorder for MockManagerClientInterface
type MockManagerClientInterfaceMockRecorder struct {
	mock *MockManagerClientInterface
}

// NewMockManagerClientInterface creates a new mock instance
func NewMockManagerClientInterface(ctrl *gomock.Controller) *MockManagerClientInterface {
	mock := &MockManagerClientInterface{ctrl: ctrl}
	mock.recorder = &MockManagerClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManagerClientInterface) EXPECT() *MockManagerClientInterfaceMockRecorder {
	return m.recorder
}

// ListVolumes mocks base method
func (m *MockManagerClientInterface) ListVolumes() ([]client.Volume, error) {
	ret := m.ctrl.Call(m, "ListVolumes")
	ret0, _ := ret[0].([]client.Volume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListVolumes indicates an expected call of ListVolumes
func (mr *MockManagerClientInterfaceMockRecorder) ListVolumes() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVolumes", reflect.TypeOf((*MockManagerClientInterface)(nil).ListVolumes))
}

// GetVolume mocks base method
func (m *MockManagerClientInterface) GetVolume(arg0 string) (*client.Volume, error) {
	ret := m.ctrl.Call(m, "GetVolume", arg0)
	ret0, _ := ret[0].(*client.Volume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVolume indicates an expected call of GetVolume
func (mr *MockManagerClientInterfaceMockRecorder) GetVolume(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVolume", reflect.TypeOf((*MockManagerClientInterface)(nil).GetVolume), arg0)
}

// ListNodes mocks base method
func (m *MockManagerClientInterface) ListNodes() ([]client.Node, error) {
	ret := m.ctrl.Call(m, "ListNodes")
	ret0, _ := ret[0].([]client.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNodes indicates an expected call of ListNodes
func (mr *MockManagerClientInterfaceMockRecorder) ListNodes() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNodes", reflect.TypeOf((*MockManagerClientInterface)(nil).ListNodes))
}

// GetNode mocks base method
func (m *MockManagerClientInterface) GetNode(arg0 string) (*client.Node, error) {
	ret := m.ctrl.Call(m, "GetNode", arg0)
	ret0, _ := ret[0].(*client.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNode indicates an expected call of GetNode
func (mr *MockManagerClientInterfaceMockRecorder) GetNode(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockManagerClientInterface)(nil).GetNode), arg0)
}

// UpdateNode mocks base method
func (m *MockManagerClientInterface) UpdateNode(arg0 *client.Node, arg1 interface{}) (*client.Node, error) {
	ret := m.ctrl.Call(m, "UpdateNode", arg0, arg1)
	ret0, _ := ret[0].(*client.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNode indicates an expected call of UpdateNode
func (mr *MockManagerClientInterfaceMockRecorder) UpdateNode(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNode", reflect.TypeOf((*MockManagerClientInterface)(nil).UpdateNode), arg0, arg1)
}

// RemoveReplica mocks base method
func (m *MockManagerClientInterface) RemoveReplica(arg0 client.Volume, arg1 string) (*client.Volume, error) {
	ret := m.ctrl.Call(m, "RemoveReplica", arg0, arg1)
	ret0, _ := ret[0].(*client.Volume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveReplica indicates an expected call of RemoveReplica
func (mr *MockManagerClientInterfaceMockRecorder) RemoveReplica(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveReplica", reflect.TypeOf((*MockManagerClientInterface)(nil).RemoveReplica), arg0, arg1)
}

// VolumeDetach mocks base method
func (m *MockManagerClientInterface) VolumeDetach(arg0 *client.Volume) (*client.Volume, error) {
	ret := m.ctrl.Call(m, "VolumeDetach", arg0)
	ret0, _ := ret[0].(*client.Volume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VolumeDetach indicates an expected call of VolumeDetach
func (mr *MockManagerClientInterfaceMockRecorder) VolumeDetach(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeDetach", reflect.TypeOf((*MockManagerClientInterface)(nil).VolumeDetach), arg0)
}

// VolumeAttach mocks base method
func (m *MockManagerClientInterface) VolumeAttach(arg0 *client.Volume, arg1 string) (*client.Volume, error) {
	ret := m.ctrl.Call(m, "VolumeAttach", arg0, arg1)
	ret0, _ := ret[0].(*client.Volume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VolumeAttach indicates an expected call of VolumeAttach
func (mr *MockManagerClientInterfaceMockRecorder) VolumeAttach(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeAttach", reflect.TypeOf((*MockManagerClientInterface)(nil).VolumeAttach), arg0, arg1)
}
