// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	api "github.com/hashicorp/consul/api"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ConsulClientInterface is an autogenerated mock type for the ConsulClientInterface type
type ConsulClientInterface struct {
	mock.Mock
}

// DeregisterService provides a mock function with given fields: ctx, serviceID
func (_m *ConsulClientInterface) DeregisterService(ctx context.Context, serviceID string) error {
	ret := _m.Called(ctx, serviceID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, serviceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLicense provides a mock function with given fields: ctx, q
func (_m *ConsulClientInterface) GetLicense(ctx context.Context, q *api.QueryOptions) (string, error) {
	ret := _m.Called(ctx, q)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *api.QueryOptions) string); ok {
		r0 = rf(ctx, q)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *api.QueryOptions) error); ok {
		r1 = rf(ctx, q)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterService provides a mock function with given fields: ctx, s
func (_m *ConsulClientInterface) RegisterService(ctx context.Context, s *api.AgentServiceRegistration) error {
	ret := _m.Called(ctx, s)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.AgentServiceRegistration) error); ok {
		r0 = rf(ctx, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
