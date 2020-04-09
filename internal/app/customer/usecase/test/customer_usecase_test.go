package test

import (
	"2020_1_drop_table/internal/app/customer/mocks"
	"2020_1_drop_table/internal/app/customer/models"
	"2020_1_drop_table/internal/app/customer/usecase"
	staffMocks "2020_1_drop_table/internal/app/staff/mocks"
	models2 "2020_1_drop_table/internal/app/staff/models"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGetPoints(t *testing.T) {
	crepo := new(mocks.Repository)
	sucase := new(staffMocks.Usecase)
	timeout := time.Second * 4
	uuid := "asdasdasd"
	retCust := models.Customer{
		CustomerID: "asdasdasd",
		CafeID:     0,
		Points:     228,
	}
	crepo.On("GetByID", mock.AnythingOfType("*context.timerCtx"), uuid).Return(retCust, nil)
	ucase := usecase.NewCustomerUsecase(crepo, sucase, timeout)

	points, err := ucase.GetPoints(context.TODO(), uuid)
	assert.Nil(t, err)
	assert.Equal(t, retCust.Points, points)
}

func TestSetPoints(t *testing.T) {
	crepo := new(mocks.Repository)
	sucase := new(staffMocks.Usecase)
	timeout := time.Second * 4
	uuid := "asdasdasd"
	newPoints := 229
	retCust := models.Customer{
		CustomerID: "asdasdasd",
		CafeID:     0,
		Points:     228,
	}

	crepo.On("GetByID", mock.AnythingOfType("*context.timerCtx"), uuid).Return(retCust, nil)
	crepo.On("SetLoyaltyPoints", mock.AnythingOfType("*context.timerCtx"), 229, uuid).Return(retCust, nil)

	returnStaff := models2.SafeStaff{
		StaffID:  228,
		Name:     "",
		Email:    "",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  false,
		CafeId:   0,
		Position: "",
	}
	sucase.On("GetFromSession", mock.AnythingOfType("*context.timerCtx")).Return(returnStaff, nil)

	ucase := usecase.NewCustomerUsecase(crepo, sucase, timeout)

	err := ucase.SetPoints(context.TODO(), uuid, newPoints)
	assert.Nil(t, err)
}
