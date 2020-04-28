package usecase

import (
	passMocks "2020_1_drop_table/internal/app/apple_passkit/mocks"
	passKitModels "2020_1_drop_table/internal/app/apple_passkit/models"
	customerMocks "2020_1_drop_table/internal/app/customer/mocks"
	customerModels "2020_1_drop_table/internal/app/customer/models"
	staffMocks "2020_1_drop_table/internal/microservices/staff/mocks"
	staffModels "2020_1_drop_table/internal/microservices/staff/models"
	"context"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCustomerUsecase_GetCustomer(t *testing.T) {
	type getCustomerTestCase struct {
		customer customerModels.Customer
		staff    staffModels.SafeStaff
		err      error
	}

	var customer customerModels.Customer
	err := faker.FakeData(&customer)
	assert.NoError(t, err)

	var staffForbidden staffModels.SafeStaff
	err = faker.FakeData(&staffForbidden)
	assert.NoError(t, err)
	staffForbidden.CafeId = customer.CafeID + 1

	staffOK := staffForbidden
	staffOK.CafeId = customer.CafeID

	testCases := []getCustomerTestCase{
		//Test OK
		{
			customer: customer,
			staff:    staffOK,
			err:      nil,
		},
		//Test forbidden
		{
			customer: customer,
			staff:    staffForbidden,
			err:      nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passMocks.Repository)
		mockCustomerRepo := new(customerMocks.Repository)
		mockStaffUcase := new(staffMocks.Usecase)

		mockStaffUcase.On("GetFromSession",
			mock.AnythingOfType("*context.timerCtx")).Return(
			testCase.staff, nil)

		customerIdMatches := func(id string) bool {
			assert.Equal(t, testCase.customer.CustomerID, id, message)
			return id == testCase.customer.CustomerID
		}
		mockCustomerRepo.On("GetByID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(customerIdMatches)).Return(
			testCase.customer, nil)

		cafeUsecase := NewCustomerUsecase(mockCustomerRepo, mockStaffUcase,
			mockPassKitRepo, time.Second*2)

		customer, err := cafeUsecase.GetCustomer(context.Background(), testCase.customer.CustomerID)

		assert.Equal(t, testCase.err, err)

		if testCase.err == nil {
			assert.Equal(t, testCase.customer, customer)
		}
	}
}

func TestCustomerUsecase_GetPoints(t *testing.T) {
	type getPointsTestCase struct {
		customer customerModels.Customer
		err      error
	}

	var customer customerModels.Customer
	err := faker.FakeData(&customer)
	assert.NoError(t, err)

	testCases := []getPointsTestCase{
		//Test OK
		{
			customer: customer,
			err:      nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passMocks.Repository)
		mockCustomerRepo := new(customerMocks.Repository)
		mockStaffUcase := new(staffMocks.Usecase)

		customerIdMatches := func(id string) bool {
			assert.Equal(t, testCase.customer.CustomerID, id, message)
			return id == testCase.customer.CustomerID
		}
		mockCustomerRepo.On("GetByID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(customerIdMatches)).Return(
			testCase.customer, nil)

		cafeUsecase := NewCustomerUsecase(mockCustomerRepo, mockStaffUcase,
			mockPassKitRepo, time.Second*2)

		customerPoints, err := cafeUsecase.GetPoints(context.Background(), testCase.customer.CustomerID)

		assert.Equal(t, testCase.err, err)

		if testCase.err == nil {
			assert.Equal(t, testCase.customer.Points, customerPoints)
		}
	}
}

func TestCustomerUsecase_SetPoints(t *testing.T) {
	type getPointsTestCase struct {
		oldCustomer customerModels.Customer
		newCustomer customerModels.Customer
		staff       staffModels.SafeStaff
		pass        passKitModels.ApplePassDB
		err         error
	}

	var oldCustomer customerModels.Customer
	err := faker.FakeData(&oldCustomer)
	assert.NoError(t, err)
	oldCustomer.Points = `{"coffee_cups": 0}`
	oldCustomer.Type = "coffee_cup"

	newCustomer := oldCustomer
	newCustomer.Points = `{"coffee_cups": 1}`

	var staffForbidden staffModels.SafeStaff
	err = faker.FakeData(&staffForbidden)
	assert.NoError(t, err)
	staffForbidden.CafeId = oldCustomer.CafeID + 1

	staffOK := staffForbidden
	staffOK.CafeId = oldCustomer.CafeID

	var pass passKitModels.ApplePassDB
	err = faker.FakeData(&staffForbidden)
	assert.NoError(t, err)
	pass.LoyaltyInfo = `{"cups_count": 10}`
	pass.Type = "coffee_cup"

	testCases := []getPointsTestCase{
		//Test OK
		{
			oldCustomer: oldCustomer,
			newCustomer: newCustomer,
			staff:       staffOK,
			pass:        pass,
			err:         nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passMocks.Repository)
		mockCustomerRepo := new(customerMocks.Repository)
		mockStaffUcase := new(staffMocks.Usecase)

		mockStaffUcase.On("GetFromSession",
			mock.AnythingOfType("*context.timerCtx")).Return(
			testCase.staff, nil)

		oldCustomerIdMatches := func(id string) bool {
			assert.Equal(t, testCase.oldCustomer.CustomerID, id, message)
			return id == testCase.oldCustomer.CustomerID
		}

		mockCustomerRepo.On("GetByID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(oldCustomerIdMatches)).Return(
			testCase.oldCustomer, nil)

		cafePassIdMatches := func(id int) bool {
			assert.Equal(t, testCase.oldCustomer.CafeID, id, message)
			return id == testCase.oldCustomer.CafeID
		}
		customerTypeMatches := func(passType string) bool {
			assert.Equal(t, testCase.oldCustomer.Type, passType, message)
			return passType == testCase.oldCustomer.Type
		}
		publishedTrue := func(published bool) bool {
			return published
		}

		mockPassKitRepo.On("GetPassByCafeID",
			mock.AnythingOfType("*context.timerCtx"),
			mock.MatchedBy(cafePassIdMatches), mock.MatchedBy(customerTypeMatches), mock.MatchedBy(publishedTrue)).Return(
			testCase.pass, nil)

		newCustomerIdMatches := func(id string) bool {
			assert.Equal(t, testCase.newCustomer.CustomerID, id, message)
			return id == testCase.newCustomer.CustomerID
		}
		newPointsMatches := func(points string) bool {
			assert.Equal(t, testCase.newCustomer.Points, points, message)
			return points == testCase.newCustomer.Points
		}

		mockCustomerRepo.On("SetLoyaltyPoints",
			mock.AnythingOfType("*context.timerCtx"),
			mock.MatchedBy(newPointsMatches), mock.MatchedBy(newCustomerIdMatches)).Return(
			testCase.newCustomer, nil)

		cafeUsecase := NewCustomerUsecase(mockCustomerRepo, mockStaffUcase,
			mockPassKitRepo, time.Second*2)

		err = cafeUsecase.SetPoints(context.Background(), testCase.newCustomer.CustomerID, testCase.newCustomer.Points)

		assert.Equal(t, testCase.err, err)
	}
}
