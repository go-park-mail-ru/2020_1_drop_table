// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	apple_pass_generator "2020_1_drop_table/internal/pkg/apple_pass_generator"
	bytes "bytes"

	mock "github.com/stretchr/testify/mock"
)

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// CreateNewPass provides a mock function with given fields: pass
func (_m *Generator) CreateNewPass(pass apple_pass_generator.ApplePass) (*bytes.Buffer, error) {
	ret := _m.Called(pass)

	var r0 *bytes.Buffer
	if rf, ok := ret.Get(0).(func(apple_pass_generator.ApplePass) *bytes.Buffer); ok {
		r0 = rf(pass)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bytes.Buffer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(apple_pass_generator.ApplePass) error); ok {
		r1 = rf(pass)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
