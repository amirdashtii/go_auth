// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"context"
	"time"

	mock "github.com/stretchr/testify/mock"
)

// NewMockInMemoryRespositoryContracts creates a new instance of InMemoryRespositoryContracts. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInMemoryRespositoryContracts(t interface {
	mock.TestingT
	Cleanup(func())
}) *InMemoryRespositoryContracts {
	mock := &InMemoryRespositoryContracts{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// InMemoryRespositoryContracts is an autogenerated mock type for the InMemoryRespositoryContracts type
type InMemoryRespositoryContracts struct {
	mock.Mock
}

type MockInMemoryRespositoryContracts_Expecter struct {
	mock *mock.Mock
}

func (_m *InMemoryRespositoryContracts) EXPECT() *MockInMemoryRespositoryContracts_Expecter {
	return &MockInMemoryRespositoryContracts_Expecter{mock: &_m.Mock}
}

// AddToken provides a mock function for the type InMemoryRespositoryContracts
func (_mock *InMemoryRespositoryContracts) AddToken(ctx context.Context, userID string, token string, expiration time.Duration) error {
	ret := _mock.Called(ctx, userID, token, expiration)

	if len(ret) == 0 {
		panic("no return value specified for AddToken")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string, time.Duration) error); ok {
		r0 = returnFunc(ctx, userID, token, expiration)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockInMemoryRespositoryContracts_AddToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddToken'
type MockInMemoryRespositoryContracts_AddToken_Call struct {
	*mock.Call
}

// AddToken is a helper method to define mock.On call
//   - ctx
//   - userID
//   - token
//   - expiration
func (_e *MockInMemoryRespositoryContracts_Expecter) AddToken(ctx interface{}, userID interface{}, token interface{}, expiration interface{}) *MockInMemoryRespositoryContracts_AddToken_Call {
	return &MockInMemoryRespositoryContracts_AddToken_Call{Call: _e.mock.On("AddToken", ctx, userID, token, expiration)}
}

func (_c *MockInMemoryRespositoryContracts_AddToken_Call) Run(run func(ctx context.Context, userID string, token string, expiration time.Duration)) *MockInMemoryRespositoryContracts_AddToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(time.Duration))
	})
	return _c
}

func (_c *MockInMemoryRespositoryContracts_AddToken_Call) Return(err error) *MockInMemoryRespositoryContracts_AddToken_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockInMemoryRespositoryContracts_AddToken_Call) RunAndReturn(run func(ctx context.Context, userID string, token string, expiration time.Duration) error) *MockInMemoryRespositoryContracts_AddToken_Call {
	_c.Call.Return(run)
	return _c
}

// FindToken provides a mock function for the type InMemoryRespositoryContracts
func (_mock *InMemoryRespositoryContracts) FindToken(ctx context.Context, userID string) (string, error) {
	ret := _mock.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for FindToken")
	}

	var r0 string
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return returnFunc(ctx, userID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = returnFunc(ctx, userID)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockInMemoryRespositoryContracts_FindToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindToken'
type MockInMemoryRespositoryContracts_FindToken_Call struct {
	*mock.Call
}

// FindToken is a helper method to define mock.On call
//   - ctx
//   - userID
func (_e *MockInMemoryRespositoryContracts_Expecter) FindToken(ctx interface{}, userID interface{}) *MockInMemoryRespositoryContracts_FindToken_Call {
	return &MockInMemoryRespositoryContracts_FindToken_Call{Call: _e.mock.On("FindToken", ctx, userID)}
}

func (_c *MockInMemoryRespositoryContracts_FindToken_Call) Run(run func(ctx context.Context, userID string)) *MockInMemoryRespositoryContracts_FindToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockInMemoryRespositoryContracts_FindToken_Call) Return(s string, err error) *MockInMemoryRespositoryContracts_FindToken_Call {
	_c.Call.Return(s, err)
	return _c
}

func (_c *MockInMemoryRespositoryContracts_FindToken_Call) RunAndReturn(run func(ctx context.Context, userID string) (string, error)) *MockInMemoryRespositoryContracts_FindToken_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveToken provides a mock function for the type InMemoryRespositoryContracts
func (_mock *InMemoryRespositoryContracts) RemoveToken(ctx context.Context, userID string) error {
	ret := _mock.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveToken")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = returnFunc(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockInMemoryRespositoryContracts_RemoveToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveToken'
type MockInMemoryRespositoryContracts_RemoveToken_Call struct {
	*mock.Call
}

// RemoveToken is a helper method to define mock.On call
//   - ctx
//   - userID
func (_e *MockInMemoryRespositoryContracts_Expecter) RemoveToken(ctx interface{}, userID interface{}) *MockInMemoryRespositoryContracts_RemoveToken_Call {
	return &MockInMemoryRespositoryContracts_RemoveToken_Call{Call: _e.mock.On("RemoveToken", ctx, userID)}
}

func (_c *MockInMemoryRespositoryContracts_RemoveToken_Call) Run(run func(ctx context.Context, userID string)) *MockInMemoryRespositoryContracts_RemoveToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockInMemoryRespositoryContracts_RemoveToken_Call) Return(err error) *MockInMemoryRespositoryContracts_RemoveToken_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockInMemoryRespositoryContracts_RemoveToken_Call) RunAndReturn(run func(ctx context.Context, userID string) error) *MockInMemoryRespositoryContracts_RemoveToken_Call {
	_c.Call.Return(run)
	return _c
}
