// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	models "2020_1_drop_table/internal/app/staff/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, st
func (_m *Repository) Add(ctx context.Context, st models.Staff) (models.Staff, error) {
	ret := _m.Called(ctx, st)

	var r0 models.Staff
	if rf, ok := ret.Get(0).(func(context.Context, models.Staff) models.Staff); ok {
		r0 = rf(ctx, st)
	} else {
		r0 = ret.Get(0).(models.Staff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Staff) error); ok {
		r1 = rf(ctx, st)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddUuid provides a mock function with given fields: ctx, uuid, id
func (_m *Repository) AddUuid(ctx context.Context, uuid string, id int) error {
	ret := _m.Called(ctx, uuid, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) error); ok {
		r0 = rf(ctx, uuid, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckIsOwner provides a mock function with given fields: ctx, staffId
func (_m *Repository) CheckIsOwner(ctx context.Context, staffId int) (bool, error) {
	ret := _m.Called(ctx, staffId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(ctx, staffId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, staffId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteStaff provides a mock function with given fields: ctx, staffId
func (_m *Repository) DeleteStaff(ctx context.Context, staffId int) error {
	ret := _m.Called(ctx, staffId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, staffId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUuid provides a mock function with given fields: ctx, uuid
func (_m *Repository) DeleteUuid(ctx context.Context, uuid string) error {
	ret := _m.Called(ctx, uuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *Repository) GetByEmail(ctx context.Context, email string) (models.Staff, error) {
	ret := _m.Called(ctx, email)

	var r0 models.Staff
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Staff); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(models.Staff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Repository) GetByID(ctx context.Context, id int) (models.Staff, error) {
	ret := _m.Called(ctx, id)

	var r0 models.Staff
	if rf, ok := ret.Get(0).(func(context.Context, int) models.Staff); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Staff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCafeId provides a mock function with given fields: ctx, uuid
func (_m *Repository) GetCafeId(ctx context.Context, uuid string) (int, error) {
	ret := _m.Called(ctx, uuid)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, uuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStaffListByOwnerId provides a mock function with given fields: ctx, ownerId
func (_m *Repository) GetStaffListByOwnerId(ctx context.Context, ownerId int) (map[string][]models.StaffByOwnerResponse, error) {
	ret := _m.Called(ctx, ownerId)

	var r0 map[string][]models.StaffByOwnerResponse
	if rf, ok := ret.Get(0).(func(context.Context, int) map[string][]models.StaffByOwnerResponse); ok {
		r0 = rf(ctx, ownerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]models.StaffByOwnerResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, ownerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, newStaff
func (_m *Repository) Update(ctx context.Context, newStaff models.SafeStaff) (models.SafeStaff, error) {
	ret := _m.Called(ctx, newStaff)

	var r0 models.SafeStaff
	if rf, ok := ret.Get(0).(func(context.Context, models.SafeStaff) models.SafeStaff); ok {
		r0 = rf(ctx, newStaff)
	} else {
		r0 = ret.Get(0).(models.SafeStaff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SafeStaff) error); ok {
		r1 = rf(ctx, newStaff)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
