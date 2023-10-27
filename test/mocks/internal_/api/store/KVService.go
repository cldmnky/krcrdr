// Code generated by mockery v2.34.1. DO NOT EDIT.

package store

import (
	context "context"

	apistore "github.com/cldmnky/krcrdr/internal/api/store"

	mock "github.com/stretchr/testify/mock"
)

// KVService is an autogenerated mock type for the KVService type
type KVService struct {
	mock.Mock
}

type KVService_Expecter struct {
	mock *mock.Mock
}

func (_m *KVService) EXPECT() *KVService_Expecter {
	return &KVService_Expecter{mock: &_m.Mock}
}

// CreateTenant provides a mock function with given fields: ctx, tenantId, tenant
func (_m *KVService) CreateTenant(ctx context.Context, tenantId string, tenant []byte) ([]byte, error) {
	ret := _m.Called(ctx, tenantId, tenant)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) ([]byte, error)); ok {
		return rf(ctx, tenantId, tenant)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) []byte); ok {
		r0 = rf(ctx, tenantId, tenant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []byte) error); ok {
		r1 = rf(ctx, tenantId, tenant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KVService_CreateTenant_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTenant'
type KVService_CreateTenant_Call struct {
	*mock.Call
}

// CreateTenant is a helper method to define mock.On call
//   - ctx context.Context
//   - tenantId string
//   - tenant []byte
func (_e *KVService_Expecter) CreateTenant(ctx interface{}, tenantId interface{}, tenant interface{}) *KVService_CreateTenant_Call {
	return &KVService_CreateTenant_Call{Call: _e.mock.On("CreateTenant", ctx, tenantId, tenant)}
}

func (_c *KVService_CreateTenant_Call) Run(run func(ctx context.Context, tenantId string, tenant []byte)) *KVService_CreateTenant_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]byte))
	})
	return _c
}

func (_c *KVService_CreateTenant_Call) Return(_a0 []byte, _a1 error) *KVService_CreateTenant_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *KVService_CreateTenant_Call) RunAndReturn(run func(context.Context, string, []byte) ([]byte, error)) *KVService_CreateTenant_Call {
	_c.Call.Return(run)
	return _c
}

// GetTenant provides a mock function with given fields: ctx, tenantId
func (_m *KVService) GetTenant(ctx context.Context, tenantId string) ([]byte, error) {
	ret := _m.Called(ctx, tenantId)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]byte, error)); ok {
		return rf(ctx, tenantId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KVService_GetTenant_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTenant'
type KVService_GetTenant_Call struct {
	*mock.Call
}

// GetTenant is a helper method to define mock.On call
//   - ctx context.Context
//   - tenantId string
func (_e *KVService_Expecter) GetTenant(ctx interface{}, tenantId interface{}) *KVService_GetTenant_Call {
	return &KVService_GetTenant_Call{Call: _e.mock.On("GetTenant", ctx, tenantId)}
}

func (_c *KVService_GetTenant_Call) Run(run func(ctx context.Context, tenantId string)) *KVService_GetTenant_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *KVService_GetTenant_Call) Return(_a0 []byte, _a1 error) *KVService_GetTenant_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *KVService_GetTenant_Call) RunAndReturn(run func(context.Context, string) ([]byte, error)) *KVService_GetTenant_Call {
	_c.Call.Return(run)
	return _c
}

// ListTenants provides a mock function with given fields: ctx
func (_m *KVService) ListTenants(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KVService_ListTenants_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListTenants'
type KVService_ListTenants_Call struct {
	*mock.Call
}

// ListTenants is a helper method to define mock.On call
//   - ctx context.Context
func (_e *KVService_Expecter) ListTenants(ctx interface{}) *KVService_ListTenants_Call {
	return &KVService_ListTenants_Call{Call: _e.mock.On("ListTenants", ctx)}
}

func (_c *KVService_ListTenants_Call) Run(run func(ctx context.Context)) *KVService_ListTenants_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *KVService_ListTenants_Call) Return(_a0 []string, _a1 error) *KVService_ListTenants_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *KVService_ListTenants_Call) RunAndReturn(run func(context.Context) ([]string, error)) *KVService_ListTenants_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx
func (_m *KVService) Watch(ctx context.Context) (<-chan apistore.KVEntry, <-chan struct{}) {
	ret := _m.Called(ctx)

	var r0 <-chan apistore.KVEntry
	var r1 <-chan struct{}
	if rf, ok := ret.Get(0).(func(context.Context) (<-chan apistore.KVEntry, <-chan struct{})); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) <-chan apistore.KVEntry); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan apistore.KVEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) <-chan struct{}); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan struct{})
		}
	}

	return r0, r1
}

// KVService_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type KVService_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
func (_e *KVService_Expecter) Watch(ctx interface{}) *KVService_Watch_Call {
	return &KVService_Watch_Call{Call: _e.mock.On("Watch", ctx)}
}

func (_c *KVService_Watch_Call) Run(run func(ctx context.Context)) *KVService_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *KVService_Watch_Call) Return(_a0 <-chan apistore.KVEntry, _a1 <-chan struct{}) *KVService_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *KVService_Watch_Call) RunAndReturn(run func(context.Context) (<-chan apistore.KVEntry, <-chan struct{})) *KVService_Watch_Call {
	_c.Call.Return(run)
	return _c
}

// NewKVService creates a new instance of KVService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKVService(t interface {
	mock.TestingT
	Cleanup(func())
}) *KVService {
	mock := &KVService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
