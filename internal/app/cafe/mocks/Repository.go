// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	models "2020_1_drop_table/internal/app/cafe/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, ca
func (_m *Repository) Add(ctx context.Context, ca models.Cafe) (models.Cafe, error) {
	ret := _m.Called(ctx, ca)

	var r0 models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, models.Cafe) models.Cafe); ok {
		r0 = rf(ctx, ca)
	} else {
		r0 = ret.Get(0).(models.Cafe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Cafe) error); ok {
		r1 = rf(ctx, ca)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Repository) GetByID(ctx context.Context, id int) (models.Cafe, error) {
	ret := _m.Called(ctx, id)

	var r0 models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, int) models.Cafe); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Cafe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByOwnerID provides a mock function with given fields: ctx, staffID
func (_m *Repository) GetByOwnerID(ctx context.Context, staffID int) ([]models.Cafe, error) {
	ret := _m.Called(ctx, staffID)

	var r0 []models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.Cafe); ok {
		r0 = rf(ctx, staffID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Cafe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, staffID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, newCafe
func (_m *Repository) Update(ctx context.Context, newCafe models.Cafe) (models.Cafe, error) {
	ret := _m.Called(ctx, newCafe)

	var r0 models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, models.Cafe) models.Cafe); ok {
		r0 = rf(ctx, newCafe)
	} else {
		r0 = ret.Get(0).(models.Cafe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Cafe) error); ok {
		r1 = rf(ctx, newCafe)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
