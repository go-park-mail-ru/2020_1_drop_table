// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "2020_1_drop_table/internal/app/customer/models"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, newCustomer
func (_m *Usecase) Add(ctx context.Context, newCustomer models.Customer) (models.Customer, error) {
	ret := _m.Called(ctx, newCustomer)

	var r0 models.Customer
	if rf, ok := ret.Get(0).(func(context.Context, models.Customer) models.Customer); ok {
		r0 = rf(ctx, newCustomer)
	} else {
		r0 = ret.Get(0).(models.Customer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Customer) error); ok {
		r1 = rf(ctx, newCustomer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCustomer provides a mock function with given fields: ctx, uuid
func (_m *Usecase) GetCustomer(ctx context.Context, uuid string) (models.Customer, error) {
	ret := _m.Called(ctx, uuid)

	var r0 models.Customer
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Customer); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Get(0).(models.Customer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, uuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPoints provides a mock function with given fields: ctx, uuid
func (_m *Usecase) GetPoints(ctx context.Context, uuid string) (string, error) {
	ret := _m.Called(ctx, uuid)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, uuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPoints provides a mock function with given fields: ctx, uuid, points
func (_m *Usecase) SetPoints(ctx context.Context, uuid string, points string) error {
	ret := _m.Called(ctx, uuid, points)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, uuid, points)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
