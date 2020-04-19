// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	models "2020_1_drop_table/internal/app/staff/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Add provides a mock function with given fields: c, newStaff
func (_m *Usecase) Add(c context.Context, newStaff models.Staff) (models.SafeStaff, error) {
	ret := _m.Called(c, newStaff)

	var r0 models.SafeStaff
	if rf, ok := ret.Get(0).(func(context.Context, models.Staff) models.SafeStaff); ok {
		r0 = rf(c, newStaff)
	} else {
		r0 = ret.Get(0).(models.SafeStaff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Staff) error); ok {
		r1 = rf(c, newStaff)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckIfStaffInOwnerCafes provides a mock function with given fields: ctx, requestUser, staffId
func (_m *Usecase) CheckIfStaffInOwnerCafes(ctx context.Context, requestUser models.SafeStaff, staffId int) (bool, error) {
	ret := _m.Called(ctx, requestUser, staffId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, models.SafeStaff, int) bool); ok {
		r0 = rf(ctx, requestUser, staffId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SafeStaff, int) error); ok {
		r1 = rf(ctx, requestUser, staffId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteQrCodes provides a mock function with given fields: uString
func (_m *Usecase) DeleteQrCodes(uString string) error {
	ret := _m.Called(uString)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(uString)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteStaffById provides a mock function with given fields: ctx, staffId
func (_m *Usecase) DeleteStaffById(ctx context.Context, staffId int) error {
	ret := _m.Called(ctx, staffId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, staffId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByEmailAndPassword provides a mock function with given fields: c, form
func (_m *Usecase) GetByEmailAndPassword(c context.Context, form models.LoginForm) (models.SafeStaff, error) {
	ret := _m.Called(c, form)

	var r0 models.SafeStaff
	if rf, ok := ret.Get(0).(func(context.Context, models.LoginForm) models.SafeStaff); ok {
		r0 = rf(c, form)
	} else {
		r0 = ret.Get(0).(models.SafeStaff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.LoginForm) error); ok {
		r1 = rf(c, form)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: c, id
func (_m *Usecase) GetByID(c context.Context, id int) (models.SafeStaff, error) {
	ret := _m.Called(c, id)

	var r0 models.SafeStaff
	if rf, ok := ret.Get(0).(func(context.Context, int) models.SafeStaff); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Get(0).(models.SafeStaff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCafeId provides a mock function with given fields: c, uuid
func (_m *Usecase) GetCafeId(c context.Context, uuid string) (int, error) {
	ret := _m.Called(c, uuid)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(c, uuid)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, uuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFromSession provides a mock function with given fields: c
func (_m *Usecase) GetFromSession(c context.Context) (models.SafeStaff, error) {
	ret := _m.Called(c)

	var r0 models.SafeStaff
	if rf, ok := ret.Get(0).(func(context.Context) models.SafeStaff); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(models.SafeStaff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetQrForStaff provides a mock function with given fields: ctx, idCafe, position
func (_m *Usecase) GetQrForStaff(ctx context.Context, idCafe int, position string) (string, error) {
	ret := _m.Called(ctx, idCafe, position)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, int, string) string); ok {
		r0 = rf(ctx, idCafe, position)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, idCafe, position)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStaffId provides a mock function with given fields: c
func (_m *Usecase) GetStaffId(c context.Context) (int, error) {
	ret := _m.Called(c)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStaffListByOwnerId provides a mock function with given fields: ctx, ownerId
func (_m *Usecase) GetStaffListByOwnerId(ctx context.Context, ownerId int) (map[string][]models.StaffByOwnerResponse, error) {
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

// IsOwner provides a mock function with given fields: c, staffId
func (_m *Usecase) IsOwner(c context.Context, staffId int) (bool, error) {
	ret := _m.Called(c, staffId)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(c, staffId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(c, staffId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: c, newStaff
func (_m *Usecase) Update(c context.Context, newStaff models.SafeStaff) (models.SafeStaff, error) {
	ret := _m.Called(c, newStaff)

	var r0 models.SafeStaff
	if rf, ok := ret.Get(0).(func(context.Context, models.SafeStaff) models.SafeStaff); ok {
		r0 = rf(c, newStaff)
	} else {
		r0 = ret.Get(0).(models.SafeStaff)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.SafeStaff) error); ok {
		r1 = rf(c, newStaff)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePosition provides a mock function with given fields: ctx, id, position
func (_m *Usecase) UpdatePosition(ctx context.Context, id int, position string) error {
	ret := _m.Called(ctx, id, position)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, id, position)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
