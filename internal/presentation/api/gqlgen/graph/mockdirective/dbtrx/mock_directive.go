// Code generated by mockery v2.10.0. DO NOT EDIT.

package dbtrx

import (
	context "context"

	graphql "github.com/99designs/gqlgen/graphql"

	mock "github.com/stretchr/testify/mock"
)

// Directive is an autogenerated mock type for the IDirective type
type Directive struct {
	mock.Mock
}

// DBTrxMiddleware provides a mock function with given fields:
func (_m *Directive) DBTrxMiddleware() func(context.Context, interface{}, graphql.Resolver) (interface{}, error) {
	ret := _m.Called()

	var r0 func(context.Context, interface{}, graphql.Resolver) (interface{}, error)
	if rf, ok := ret.Get(0).(func() func(context.Context, interface{}, graphql.Resolver) (interface{}, error)); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func(context.Context, interface{}, graphql.Resolver) (interface{}, error))
		}
	}

	return r0
}
