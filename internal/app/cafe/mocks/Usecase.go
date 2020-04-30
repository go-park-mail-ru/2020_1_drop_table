// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	models "2020_1_drop_table/internal/app/cafe/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Add provides a mock function with given fields: c, newCafe
func (_m *Usecase) Add(c context.Context, newCafe models.Cafe) (models.Cafe, error) {
	ret := _m.Called(c, newCafe)

	var r0 models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, models.Cafe) models.Cafe); ok {
		r0 = rf(c, newCafe)
	} else {
		r0 = ret.Get(0).(models.Cafe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Cafe) error); ok {
		r1 = rf(c, newCafe)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllCafes provides a mock function with given fields: ctx, since, limit
func (_m *Usecase) GetAllCafes(ctx context.Context, since int, limit int) ([]models.Cafe, error) {
	ret := _m.Called(ctx, since, limit)

	var r0 []models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []models.Cafe); ok {
		r0 = rf(ctx, since, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Cafe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, since, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: c, id
func (_m *Usecase) GetByID(c context.Context, id int) (models.Cafe, error) {
	ret := _m.Called(c, id)

	var r0 models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, int) models.Cafe); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Get(0).(models.Cafe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByOwnerID provides a mock function with given fields: c
func (_m *Usecase) GetByOwnerID(c context.Context) ([]models.Cafe, error) {
	ret := _m.Called(c)

	var r0 []models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context) []models.Cafe); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Cafe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: c, newCafe
func (_m *Usecase) Update(c context.Context, newCafe models.Cafe) (models.Cafe, error) {
	ret := _m.Called(c, newCafe)

	var r0 models.Cafe
	if rf, ok := ret.Get(0).(func(context.Context, models.Cafe) models.Cafe); ok {
		r0 = rf(c, newCafe)
	} else {
		r0 = ret.Get(0).(models.Cafe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Cafe) error); ok {
		r1 = rf(c, newCafe)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
