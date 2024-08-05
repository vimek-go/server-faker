// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	api "github.com/vimek-go/server-faker/internal/pkg/api"
	dto "github.com/vimek-go/server-faker/internal/pkg/parser/dto"

	mock "github.com/stretchr/testify/mock"
)

// FactoryMock is an autogenerated mock type for the Factory type
type FactoryMock struct {
	mock.Mock
}

// CreateEndpoint provides a mock function with given fields: endpoint, baseDir
func (_m *FactoryMock) CreateEndpoint(endpoint dto.Endpoint, baseDir string) (api.Handler, error) {
	ret := _m.Called(endpoint, baseDir)

	if len(ret) == 0 {
		panic("no return value specified for CreateEndpoint")
	}

	var r0 api.Handler
	var r1 error
	if rf, ok := ret.Get(0).(func(dto.Endpoint, string) (api.Handler, error)); ok {
		return rf(endpoint, baseDir)
	}
	if rf, ok := ret.Get(0).(func(dto.Endpoint, string) api.Handler); ok {
		r0 = rf(endpoint, baseDir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.Handler)
		}
	}

	if rf, ok := ret.Get(1).(func(dto.Endpoint, string) error); ok {
		r1 = rf(endpoint, baseDir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateProxyEndpoint provides a mock function with given fields: endpoint
func (_m *FactoryMock) CreateProxyEndpoint(endpoint dto.Endpoint) (api.Handler, error) {
	ret := _m.Called(endpoint)

	if len(ret) == 0 {
		panic("no return value specified for CreateProxyEndpoint")
	}

	var r0 api.Handler
	var r1 error
	if rf, ok := ret.Get(0).(func(dto.Endpoint) (api.Handler, error)); ok {
		return rf(endpoint)
	}
	if rf, ok := ret.Get(0).(func(dto.Endpoint) api.Handler); ok {
		r0 = rf(endpoint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.Handler)
		}
	}

	if rf, ok := ret.Get(1).(func(dto.Endpoint) error); ok {
		r1 = rf(endpoint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateResponseEndpoint provides a mock function with given fields: endpoint, baseDir
func (_m *FactoryMock) CreateResponseEndpoint(endpoint dto.Endpoint, baseDir string) (api.ResponseHandler, error) {
	ret := _m.Called(endpoint, baseDir)

	if len(ret) == 0 {
		panic("no return value specified for CreateResponseEndpoint")
	}

	var r0 api.ResponseHandler
	var r1 error
	if rf, ok := ret.Get(0).(func(dto.Endpoint, string) (api.ResponseHandler, error)); ok {
		return rf(endpoint, baseDir)
	}
	if rf, ok := ret.Get(0).(func(dto.Endpoint, string) api.ResponseHandler); ok {
		r0 = rf(endpoint, baseDir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.ResponseHandler)
		}
	}

	if rf, ok := ret.Get(1).(func(dto.Endpoint, string) error); ok {
		r1 = rf(endpoint, baseDir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFactoryMock creates a new instance of FactoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFactoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *FactoryMock {
	mock := &FactoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}