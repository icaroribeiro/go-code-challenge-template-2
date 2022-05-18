// Code generated by mockery v2.10.0. DO NOT EDIT.

package mockvalidator

import mock "github.com/stretchr/testify/mock"

// Validator is an autogenerated mock type for the IValidator type
type Validator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: i
func (_m *Validator) Validate(i interface{}) error {
	ret := _m.Called(i)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateWithTags provides a mock function with given fields: i, tags
func (_m *Validator) ValidateWithTags(i interface{}, tags string) error {
	ret := _m.Called(i, tags)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string) error); ok {
		r0 = rf(i, tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}