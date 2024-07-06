// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	notify "github.com/trb1maker/gobannercircle/internal/notify"
)

// Notifier is an autogenerated mock type for the Notifier type
type Notifier struct {
	mock.Mock
}

// Notify provides a mock function with given fields: ctx, message
func (_m *Notifier) Notify(ctx context.Context, message notify.Message) error {
	ret := _m.Called(ctx, message)

	if len(ret) == 0 {
		panic("no return value specified for Notify")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, notify.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewNotifier creates a new instance of Notifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *Notifier {
	mock := &Notifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
