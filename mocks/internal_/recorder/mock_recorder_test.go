// Code generated by mockery v2.34.1. DO NOT EDIT.

package recorder

import (
	mock "github.com/stretchr/testify/mock"
	v1 "k8s.io/api/admission/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Recorder is an autogenerated mock type for the Recorder type
type Recorder struct {
	mock.Mock
}

type Recorder_Expecter struct {
	mock *mock.Mock
}

func (_m *Recorder) EXPECT() *Recorder_Expecter {
	return &Recorder_Expecter{mock: &_m.Mock}
}

// FromAdmissionRequest provides a mock function with given fields: oldObject, newObject, req
func (_m *Recorder) FromAdmissionRequest(oldObject *unstructured.Unstructured, newObject *unstructured.Unstructured, req *v1.AdmissionRequest) error {
	ret := _m.Called(oldObject, newObject, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(*unstructured.Unstructured, *unstructured.Unstructured, *v1.AdmissionRequest) error); ok {
		r0 = rf(oldObject, newObject, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Recorder_FromAdmissionRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FromAdmissionRequest'
type Recorder_FromAdmissionRequest_Call struct {
	*mock.Call
}

// FromAdmissionRequest is a helper method to define mock.On call
//   - oldObject *unstructured.Unstructured
//   - newObject *unstructured.Unstructured
//   - req *v1.AdmissionRequest
func (_e *Recorder_Expecter) FromAdmissionRequest(oldObject interface{}, newObject interface{}, req interface{}) *Recorder_FromAdmissionRequest_Call {
	return &Recorder_FromAdmissionRequest_Call{Call: _e.mock.On("FromAdmissionRequest", oldObject, newObject, req)}
}

func (_c *Recorder_FromAdmissionRequest_Call) Run(run func(oldObject *unstructured.Unstructured, newObject *unstructured.Unstructured, req *v1.AdmissionRequest)) *Recorder_FromAdmissionRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*unstructured.Unstructured), args[1].(*unstructured.Unstructured), args[2].(*v1.AdmissionRequest))
	})
	return _c
}

func (_c *Recorder_FromAdmissionRequest_Call) Return(_a0 error) *Recorder_FromAdmissionRequest_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Recorder_FromAdmissionRequest_Call) RunAndReturn(run func(*unstructured.Unstructured, *unstructured.Unstructured, *v1.AdmissionRequest) error) *Recorder_FromAdmissionRequest_Call {
	_c.Call.Return(run)
	return _c
}

// OperationType provides a mock function with given fields:
func (_m *Recorder) OperationType() v1.Operation {
	ret := _m.Called()

	var r0 v1.Operation
	if rf, ok := ret.Get(0).(func() v1.Operation); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(v1.Operation)
	}

	return r0
}

// Recorder_OperationType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OperationType'
type Recorder_OperationType_Call struct {
	*mock.Call
}

// OperationType is a helper method to define mock.On call
func (_e *Recorder_Expecter) OperationType() *Recorder_OperationType_Call {
	return &Recorder_OperationType_Call{Call: _e.mock.On("OperationType")}
}

func (_c *Recorder_OperationType_Call) Run(run func()) *Recorder_OperationType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Recorder_OperationType_Call) Return(_a0 v1.Operation) *Recorder_OperationType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Recorder_OperationType_Call) RunAndReturn(run func() v1.Operation) *Recorder_OperationType_Call {
	_c.Call.Return(run)
	return _c
}

// SendToApiServer provides a mock function with given fields:
func (_m *Recorder) SendToApiServer() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Recorder_SendToApiServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendToApiServer'
type Recorder_SendToApiServer_Call struct {
	*mock.Call
}

// SendToApiServer is a helper method to define mock.On call
func (_e *Recorder_Expecter) SendToApiServer() *Recorder_SendToApiServer_Call {
	return &Recorder_SendToApiServer_Call{Call: _e.mock.On("SendToApiServer")}
}

func (_c *Recorder_SendToApiServer_Call) Run(run func()) *Recorder_SendToApiServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Recorder_SendToApiServer_Call) Return(_a0 error) *Recorder_SendToApiServer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Recorder_SendToApiServer_Call) RunAndReturn(run func() error) *Recorder_SendToApiServer_Call {
	_c.Call.Return(run)
	return _c
}

// ToYaml provides a mock function with given fields:
func (_m *Recorder) ToYaml() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Recorder_ToYaml_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ToYaml'
type Recorder_ToYaml_Call struct {
	*mock.Call
}

// ToYaml is a helper method to define mock.On call
func (_e *Recorder_Expecter) ToYaml() *Recorder_ToYaml_Call {
	return &Recorder_ToYaml_Call{Call: _e.mock.On("ToYaml")}
}

func (_c *Recorder_ToYaml_Call) Run(run func()) *Recorder_ToYaml_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Recorder_ToYaml_Call) Return(_a0 string, _a1 error) *Recorder_ToYaml_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Recorder_ToYaml_Call) RunAndReturn(run func() (string, error)) *Recorder_ToYaml_Call {
	_c.Call.Return(run)
	return _c
}

// NewRecorder creates a new instance of Recorder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRecorder(t interface {
	mock.TestingT
	Cleanup(func())
}) *Recorder {
	mock := &Recorder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
