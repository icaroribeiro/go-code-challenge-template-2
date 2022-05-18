// Code generated by mockery v2.10.0. DO NOT EDIT.

package user

import (
	model "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/domain/model"
	user "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/core/ports/infrastructure/storage/datastore/repository/user"
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// Repository is an autogenerated mock type for the IRepository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *Repository) Create(_a0 model.User) (model.User, error) {
	ret := _m.Called(_a0)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(model.User) model.User); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *Repository) GetAll() (model.Users, error) {
	ret := _m.Called()

	var r0 model.Users
	if rf, ok := ret.Get(0).(func() model.Users); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Users)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithDBTrx provides a mock function with given fields: dbTrx
func (_m *Repository) WithDBTrx(dbTrx *gorm.DB) user.IRepository {
	ret := _m.Called(dbTrx)

	var r0 user.IRepository
	if rf, ok := ret.Get(0).(func(*gorm.DB) user.IRepository); ok {
		r0 = rf(dbTrx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(user.IRepository)
		}
	}

	return r0
}