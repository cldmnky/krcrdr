// Code generated by mockery v2.34.1. DO NOT EDIT.

package options

import (
	cobra "github.com/spf13/cobra"
	mock "github.com/stretchr/testify/mock"
)

// Interface is an autogenerated mock type for the Interface type
type Interface struct {
	mock.Mock
}

type Interface_Expecter struct {
	mock *mock.Mock
}

func (_m *Interface) EXPECT() *Interface_Expecter {
	return &Interface_Expecter{mock: &_m.Mock}
}

// AddFlags provides a mock function with given fields: cmd
func (_m *Interface) AddFlags(cmd *cobra.Command) {
	_m.Called(cmd)
}

// Interface_AddFlags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFlags'
type Interface_AddFlags_Call struct {
	*mock.Call
}

// AddFlags is a helper method to define mock.On call
//   - cmd *cobra.Command
func (_e *Interface_Expecter) AddFlags(cmd interface{}) *Interface_AddFlags_Call {
	return &Interface_AddFlags_Call{Call: _e.mock.On("AddFlags", cmd)}
}

func (_c *Interface_AddFlags_Call) Run(run func(cmd *cobra.Command)) *Interface_AddFlags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*cobra.Command))
	})
	return _c
}

func (_c *Interface_AddFlags_Call) Return() *Interface_AddFlags_Call {
	_c.Call.Return()
	return _c
}

func (_c *Interface_AddFlags_Call) RunAndReturn(run func(*cobra.Command)) *Interface_AddFlags_Call {
	_c.Call.Return(run)
	return _c
}

// NewInterface creates a new instance of Interface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *Interface {
	mock := &Interface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
