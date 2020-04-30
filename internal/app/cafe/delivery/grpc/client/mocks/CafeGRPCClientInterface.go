// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	models "2020_1_drop_table/internal/app/cafe/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CafeGRPCClientInterface is an autogenerated mock type for the CafeGRPCClientInterface type
type CafeGRPCClientInterface struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *CafeGRPCClientInterface) GetByID(ctx context.Context, id int) (models.Cafe, error) {
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