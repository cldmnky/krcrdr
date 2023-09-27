// Code generated by mockery v2.34.1. DO NOT EDIT.

package record

import (
	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// MiddlewareFunc is an autogenerated mock type for the MiddlewareFunc type
type MiddlewareFunc struct {
	mock.Mock
}

type MiddlewareFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *MiddlewareFunc) EXPECT() *MiddlewareFunc_Expecter {
	return &MiddlewareFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: c
func (_m *MiddlewareFunc) Execute(c *gin.Context) {
	_m.Called(c)
}

// MiddlewareFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MiddlewareFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - c *gin.Context
func (_e *MiddlewareFunc_Expecter) Execute(c interface{}) *MiddlewareFunc_Execute_Call {
	return &MiddlewareFunc_Execute_Call{Call: _e.mock.On("Execute", c)}
}

func (_c *MiddlewareFunc_Execute_Call) Run(run func(c *gin.Context)) *MiddlewareFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gin.Context))
	})
	return _c
}

func (_c *MiddlewareFunc_Execute_Call) Return() *MiddlewareFunc_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MiddlewareFunc_Execute_Call) RunAndReturn(run func(*gin.Context)) *MiddlewareFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMiddlewareFunc creates a new instance of MiddlewareFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMiddlewareFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MiddlewareFunc {
	mock := &MiddlewareFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
