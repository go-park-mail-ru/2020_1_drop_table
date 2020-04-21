// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "2020_1_drop_table/internal/app/customer/models"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, cu
func (_m *Repository) Add(ctx context.Context, cu models.Customer) (models.Customer, error) {
	ret := _m.Called(ctx, cu)

	var r0 models.Customer
	if rf, ok := ret.Get(0).(func(context.Context, models.Customer) models.Customer); ok {
		r0 = rf(ctx, cu)
	} else {
		r0 = ret.Get(0).(models.Customer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Customer) error); ok {
		r1 = rf(ctx, cu)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByID provides a mock function with given fields: ctx, customerID
func (_m *Repository) DeleteByID(ctx context.Context, customerID string) error {
	ret := _m.Called(ctx, customerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, customerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, customerID
func (_m *Repository) GetByID(ctx context.Context, customerID string) (models.Customer, error) {
	ret := _m.Called(ctx, customerID)

	var r0 models.Customer
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Customer); ok {
		r0 = rf(ctx, customerID)
	} else {
		r0 = ret.Get(0).(models.Customer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncrementSum provides a mock function with given fields: ctx, sum, uuid
func (_m *Repository) IncrementSum(ctx context.Context, sum float32, uuid string) error {
	ret := _m.Called(ctx, sum, uuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, float32, string) error); ok {
		r0 = rf(ctx, sum, uuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetLoyaltyPoints provides a mock function with given fields: ctx, points, customerID
func (_m *Repository) SetLoyaltyPoints(ctx context.Context, points string, customerID string) (models.Customer, error) {
	ret := _m.Called(ctx, points, customerID)

	var r0 models.Customer
	if rf, ok := ret.Get(0).(func(context.Context, int, string) models.Customer); ok {
		r0 = rf(ctx, points, customerID)
	} else {
		r0 = ret.Get(0).(models.Customer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, points, customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
