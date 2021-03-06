// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	models "2020_1_drop_table/internal/app/cafe/models"
	context "context"

	mock "github.com/stretchr/testify/mock"

	statisticsmodels "2020_1_drop_table/internal/app/statistics/models"

	time "time"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// AddData provides a mock function with given fields: jsonData, _a1, clientUUID, staffID, cafeId
func (_m *Repository) AddData(jsonData string, _a1 time.Time, clientUUID string, staffID int, cafeId int) error {
	ret := _m.Called(jsonData, _a1, clientUUID, staffID, cafeId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Time, string, int, int) error); ok {
		r0 = rf(jsonData, _a1, clientUUID, staffID, cafeId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetGraphsDataFromRepo provides a mock function with given fields: ctx, cafeList, typ, since, to
func (_m *Repository) GetGraphsDataFromRepo(ctx context.Context, cafeList []models.Cafe, typ string, since string, to string) ([]statisticsmodels.StatisticsGraphRawStruct, error) {
	ret := _m.Called(ctx, cafeList, typ, since, to)

	var r0 []statisticsmodels.StatisticsGraphRawStruct
	if rf, ok := ret.Get(0).(func(context.Context, []models.Cafe, string, string, string) []statisticsmodels.StatisticsGraphRawStruct); ok {
		r0 = rf(ctx, cafeList, typ, since, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]statisticsmodels.StatisticsGraphRawStruct)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []models.Cafe, string, string, string) error); ok {
		r1 = rf(ctx, cafeList, typ, since, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWorkerDataFromRepo provides a mock function with given fields: ctx, staffId, limit, since
func (_m *Repository) GetWorkerDataFromRepo(ctx context.Context, staffId int, limit int, since int) ([]statisticsmodels.StatisticsStruct, error) {
	ret := _m.Called(ctx, staffId, limit, since)

	var r0 []statisticsmodels.StatisticsStruct
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) []statisticsmodels.StatisticsStruct); ok {
		r0 = rf(ctx, staffId, limit, since)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]statisticsmodels.StatisticsStruct)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int, int) error); ok {
		r1 = rf(ctx, staffId, limit, since)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
