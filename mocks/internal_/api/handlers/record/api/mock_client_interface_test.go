// Code generated by mockery v2.34.1. DO NOT EDIT.

package api

import (
	context "context"
	http "net/http"

	io "io"

	mock "github.com/stretchr/testify/mock"

	recordapi "github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
)

// ClientInterface is an autogenerated mock type for the ClientInterface type
type ClientInterface struct {
	mock.Mock
}

type ClientInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *ClientInterface) EXPECT() *ClientInterface_Expecter {
	return &ClientInterface_Expecter{mock: &_m.Mock}
}

// AddRecord provides a mock function with given fields: ctx, body, reqEditors
func (_m *ClientInterface) AddRecord(ctx context.Context, body recordapi.Record, reqEditors ...recordapi.RequestEditorFn) (*http.Response, error) {
	_va := make([]interface{}, len(reqEditors))
	for _i := range reqEditors {
		_va[_i] = reqEditors[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, body)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, recordapi.Record, ...recordapi.RequestEditorFn) (*http.Response, error)); ok {
		return rf(ctx, body, reqEditors...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, recordapi.Record, ...recordapi.RequestEditorFn) *http.Response); ok {
		r0 = rf(ctx, body, reqEditors...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, recordapi.Record, ...recordapi.RequestEditorFn) error); ok {
		r1 = rf(ctx, body, reqEditors...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientInterface_AddRecord_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddRecord'
type ClientInterface_AddRecord_Call struct {
	*mock.Call
}

// AddRecord is a helper method to define mock.On call
//   - ctx context.Context
//   - body recordapi.Record
//   - reqEditors ...recordapi.RequestEditorFn
func (_e *ClientInterface_Expecter) AddRecord(ctx interface{}, body interface{}, reqEditors ...interface{}) *ClientInterface_AddRecord_Call {
	return &ClientInterface_AddRecord_Call{Call: _e.mock.On("AddRecord",
		append([]interface{}{ctx, body}, reqEditors...)...)}
}

func (_c *ClientInterface_AddRecord_Call) Run(run func(ctx context.Context, body recordapi.Record, reqEditors ...recordapi.RequestEditorFn)) *ClientInterface_AddRecord_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]recordapi.RequestEditorFn, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(recordapi.RequestEditorFn)
			}
		}
		run(args[0].(context.Context), args[1].(recordapi.Record), variadicArgs...)
	})
	return _c
}

func (_c *ClientInterface_AddRecord_Call) Return(_a0 *http.Response, _a1 error) *ClientInterface_AddRecord_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ClientInterface_AddRecord_Call) RunAndReturn(run func(context.Context, recordapi.Record, ...recordapi.RequestEditorFn) (*http.Response, error)) *ClientInterface_AddRecord_Call {
	_c.Call.Return(run)
	return _c
}

// AddRecordWithBody provides a mock function with given fields: ctx, contentType, body, reqEditors
func (_m *ClientInterface) AddRecordWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...recordapi.RequestEditorFn) (*http.Response, error) {
	_va := make([]interface{}, len(reqEditors))
	for _i := range reqEditors {
		_va[_i] = reqEditors[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, contentType, body)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, io.Reader, ...recordapi.RequestEditorFn) (*http.Response, error)); ok {
		return rf(ctx, contentType, body, reqEditors...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, io.Reader, ...recordapi.RequestEditorFn) *http.Response); ok {
		r0 = rf(ctx, contentType, body, reqEditors...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, io.Reader, ...recordapi.RequestEditorFn) error); ok {
		r1 = rf(ctx, contentType, body, reqEditors...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientInterface_AddRecordWithBody_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddRecordWithBody'
type ClientInterface_AddRecordWithBody_Call struct {
	*mock.Call
}

// AddRecordWithBody is a helper method to define mock.On call
//   - ctx context.Context
//   - contentType string
//   - body io.Reader
//   - reqEditors ...recordapi.RequestEditorFn
func (_e *ClientInterface_Expecter) AddRecordWithBody(ctx interface{}, contentType interface{}, body interface{}, reqEditors ...interface{}) *ClientInterface_AddRecordWithBody_Call {
	return &ClientInterface_AddRecordWithBody_Call{Call: _e.mock.On("AddRecordWithBody",
		append([]interface{}{ctx, contentType, body}, reqEditors...)...)}
}

func (_c *ClientInterface_AddRecordWithBody_Call) Run(run func(ctx context.Context, contentType string, body io.Reader, reqEditors ...recordapi.RequestEditorFn)) *ClientInterface_AddRecordWithBody_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]recordapi.RequestEditorFn, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(recordapi.RequestEditorFn)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(io.Reader), variadicArgs...)
	})
	return _c
}

func (_c *ClientInterface_AddRecordWithBody_Call) Return(_a0 *http.Response, _a1 error) *ClientInterface_AddRecordWithBody_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ClientInterface_AddRecordWithBody_Call) RunAndReturn(run func(context.Context, string, io.Reader, ...recordapi.RequestEditorFn) (*http.Response, error)) *ClientInterface_AddRecordWithBody_Call {
	_c.Call.Return(run)
	return _c
}

// ListRecords provides a mock function with given fields: ctx, reqEditors
func (_m *ClientInterface) ListRecords(ctx context.Context, reqEditors ...recordapi.RequestEditorFn) (*http.Response, error) {
	_va := make([]interface{}, len(reqEditors))
	for _i := range reqEditors {
		_va[_i] = reqEditors[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...recordapi.RequestEditorFn) (*http.Response, error)); ok {
		return rf(ctx, reqEditors...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...recordapi.RequestEditorFn) *http.Response); ok {
		r0 = rf(ctx, reqEditors...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...recordapi.RequestEditorFn) error); ok {
		r1 = rf(ctx, reqEditors...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientInterface_ListRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListRecords'
type ClientInterface_ListRecords_Call struct {
	*mock.Call
}

// ListRecords is a helper method to define mock.On call
//   - ctx context.Context
//   - reqEditors ...recordapi.RequestEditorFn
func (_e *ClientInterface_Expecter) ListRecords(ctx interface{}, reqEditors ...interface{}) *ClientInterface_ListRecords_Call {
	return &ClientInterface_ListRecords_Call{Call: _e.mock.On("ListRecords",
		append([]interface{}{ctx}, reqEditors...)...)}
}

func (_c *ClientInterface_ListRecords_Call) Run(run func(ctx context.Context, reqEditors ...recordapi.RequestEditorFn)) *ClientInterface_ListRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]recordapi.RequestEditorFn, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(recordapi.RequestEditorFn)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *ClientInterface_ListRecords_Call) Return(_a0 *http.Response, _a1 error) *ClientInterface_ListRecords_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ClientInterface_ListRecords_Call) RunAndReturn(run func(context.Context, ...recordapi.RequestEditorFn) (*http.Response, error)) *ClientInterface_ListRecords_Call {
	_c.Call.Return(run)
	return _c
}

// NewClientInterface creates a new instance of ClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClientInterface {
	mock := &ClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
