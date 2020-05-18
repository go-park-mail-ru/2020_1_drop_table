package usecase_test

import (
	cafeMocks "2020_1_drop_table/internal/app/cafe/mocks"
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	_cafeUsecase "2020_1_drop_table/internal/app/cafe/usecase"
	globalModels "2020_1_drop_table/internal/app/models"
	staffClientMock "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client/mocks"
	staffModels "2020_1_drop_table/internal/microservices/staff/models"
	geoMocks "2020_1_drop_table/internal/pkg/google_geocoder/mocks"
	"context"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	type addTestCase struct {
		inputCafe    cafeModels.Cafe
		expectedCafe cafeModels.Cafe
		staff        staffModels.SafeStaff
		err          error
	}

	mockCafeRepo := new(cafeMocks.Repository)
	mockStaffGRPCClient := new(staffClientMock.StaffClientInterface)
	mockGeoCoder := new(geoMocks.GoogleGeoCoder)

	cafeUsecase := _cafeUsecase.NewCafeUsecase(mockCafeRepo, mockStaffGRPCClient, time.Second*2, mockGeoCoder)
	var owner staffModels.SafeStaff
	err := faker.FakeData(&owner)
	assert.NoError(t, err)
	owner.IsOwner = true

	var inputCafe cafeModels.Cafe
	err = faker.FakeData(&inputCafe)
	assert.NoError(t, err)
	inputCafe.StaffID = owner.StaffID
	expectedCafe := inputCafe
	inputCafe.StaffID = 0

	var notOwner staffModels.SafeStaff
	notOwner.IsOwner = false

	var anonymous staffModels.SafeStaff
	anonymous.StaffID = -1

	testCases := []addTestCase{
		//Test OK
		{
			inputCafe:    inputCafe,
			expectedCafe: expectedCafe,
			staff:        owner,
			err:          nil,
		},
		//Test not owner user
		{
			inputCafe:    cafeModels.Cafe{},
			expectedCafe: cafeModels.Cafe{},
			staff:        notOwner,
			err:          globalModels.ErrForbidden,
		},
		//Test anonymous user
		{
			inputCafe:    inputCafe,
			expectedCafe: expectedCafe,
			staff:        anonymous,
			err:          globalModels.ErrForbidden,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockStaffGRPCClient.On("GetById",
			mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("int")).Return(
			testCase.staff, nil)

		mockGeoCoder.On("GetGeoByAddress",
			mock.AnythingOfType("string")).Return(testCase.inputCafe.Address, nil)

		cafeMatches := func(c cafeModels.Cafe) bool {
			assert.Equal(t, testCase.expectedCafe, c, message)
			return c == testCase.expectedCafe
		}

		mockCafeRepo.On("Add",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(cafeMatches)).Return(
			testCase.expectedCafe, nil)

		sessionUserID := testCase.staff.StaffID
		session := sessions.Session{Values: map[interface{}]interface{}{"userID": sessionUserID}}
		c := context.WithValue(context.Background(), "session", &session)

		realCafe, err := cafeUsecase.Add(c, testCase.inputCafe)
		if err == nil {
			assert.NoError(t, err, message)
			assert.Equal(t, testCase.expectedCafe, realCafe, message)
		} else {
			assert.Equal(t, testCase.err, err, message)
		}
	}
}

func TestGetByOwnerID(t *testing.T) {
	type getByOwnerTestCase struct {
		staffID       int
		expectedCafes []cafeModels.Cafe
		err           error
	}

	mockCafeRepo := new(cafeMocks.Repository)
	mockStaffGRPCClient := new(staffClientMock.StaffClientInterface)
	mockGeoCoder := new(geoMocks.GoogleGeoCoder)
	cafeUsecase := _cafeUsecase.NewCafeUsecase(mockCafeRepo, mockStaffGRPCClient, time.Second*2, mockGeoCoder)

	var owner staffModels.SafeStaff
	err := faker.FakeData(&owner)
	assert.NoError(t, err)
	owner.IsOwner = true

	cafeArray := make([]cafeModels.Cafe, 4, 4)
	err = faker.FakeData(&cafeArray)
	assert.NoError(t, err)
	for _, cafe := range cafeArray {
		cafe.StaffID = owner.StaffID
	}

	var anonymous staffModels.SafeStaff
	anonymous.StaffID = -1

	testCases := []getByOwnerTestCase{
		//Test OK
		{
			staffID:       5,
			expectedCafes: cafeArray,
			err:           nil,
		},
		//Test anonymous user
		{
			staffID:       -1,
			expectedCafes: make([]cafeModels.Cafe, 0, 0),
			err:           globalModels.ErrForbidden,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		idMatches := func(id int) bool {
			assert.Equal(t, testCase.staffID, id, message)
			return testCase.staffID == id
		}

		mockCafeRepo.On("GetByOwnerID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(idMatches)).Return(
			testCase.expectedCafes, testCase.err)

		session := sessions.Session{Values: map[interface{}]interface{}{"userID": testCase.staffID}}
		c := context.WithValue(context.Background(), "session", &session)

		realCafes, err := cafeUsecase.GetByOwnerID(c)
		assert.Equal(t, testCase.expectedCafes, realCafes, message)
		assert.Equal(t, testCase.err, err, message)
	}
}

func TestGetByID(t *testing.T) {
	type getByOwnerTestCase struct {
		staffID      int
		expectedCafe cafeModels.Cafe
		err          error
	}

	mockCafeRepo := new(cafeMocks.Repository)
	mockStaffGRPCClient := new(staffClientMock.StaffClientInterface)
	mockGeoCoder := new(geoMocks.GoogleGeoCoder)

	cafeUsecase := _cafeUsecase.NewCafeUsecase(mockCafeRepo, mockStaffGRPCClient, time.Second*2, mockGeoCoder)

	ownerID := 1

	var expectedCafe cafeModels.Cafe
	err := faker.FakeData(&expectedCafe)
	assert.NoError(t, err)
	expectedCafe.StaffID = ownerID

	testCases := []getByOwnerTestCase{
		//Test OK
		{
			staffID:      ownerID,
			expectedCafe: expectedCafe,
			err:          nil,
		},
		//Test OK (anonymous user)
		{
			staffID:      -1,
			expectedCafe: expectedCafe,
			err:          nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		idMatches := func(id int) bool {
			assert.Equal(t, testCase.expectedCafe.CafeID, id, message)
			return testCase.expectedCafe.CafeID == id
		}

		mockCafeRepo.On("GetByID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(idMatches)).Return(
			testCase.expectedCafe, testCase.err)

		realCafee, err := cafeUsecase.GetByID(context.Background(), testCase.expectedCafe.CafeID)
		assert.Equal(t, testCase.expectedCafe, realCafee, message)
		assert.Equal(t, testCase.err, err, message)

	}
}

func TestUpdate(t *testing.T) {
	type updateTestCase struct {
		staff      staffModels.SafeStaff
		oldCafe    cafeModels.Cafe
		inputCafe  cafeModels.Cafe
		outputCafe cafeModels.Cafe
		err        error
	}

	mockCafeRepo := new(cafeMocks.Repository)
	mockStaffGRPCClient := new(staffClientMock.StaffClientInterface)
	mockGeoCoder := new(geoMocks.GoogleGeoCoder)

	cafeUsecase := _cafeUsecase.NewCafeUsecase(mockCafeRepo, mockStaffGRPCClient, time.Second*2, mockGeoCoder)

	var owner staffModels.SafeStaff
	err := faker.FakeData(&owner)
	assert.NoError(t, err)
	owner.IsOwner = true

	var notOwner staffModels.SafeStaff
	notOwner.IsOwner = false

	var anonymous staffModels.SafeStaff
	anonymous.StaffID = -1

	var oldCafe cafeModels.Cafe
	err = faker.FakeData(&oldCafe)
	assert.NoError(t, err)
	oldCafe.StaffID = owner.StaffID

	inputCafe := oldCafe
	inputCafe.StaffID = 0
	inputCafe.CafeName = "NEW CAFE NAME"

	outputCafe := inputCafe
	outputCafe.StaffID = owner.StaffID

	testCases := []updateTestCase{
		//Test OK
		{
			staff:      owner,
			oldCafe:    oldCafe,
			inputCafe:  inputCafe,
			outputCafe: outputCafe,
			err:        nil,
		},
		//Test not owner user
		{
			staff:      notOwner,
			oldCafe:    oldCafe,
			inputCafe:  inputCafe,
			outputCafe: outputCafe,
			err:        globalModels.ErrForbidden,
		},
		//Test anonymous user
		{
			staff:      anonymous,
			oldCafe:    oldCafe,
			inputCafe:  inputCafe,
			outputCafe: outputCafe,
			err:        globalModels.ErrForbidden,
		},
		//Test someone else's cafe
		{
			staff:      anonymous,
			oldCafe:    oldCafe,
			inputCafe:  inputCafe,
			outputCafe: outputCafe,
			err:        globalModels.ErrForbidden,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		cafeIdMatches := func(id int) bool {
			assert.Equal(t, testCase.oldCafe.CafeID, id, message)
			return id == testCase.oldCafe.CafeID
		}

		mockCafeRepo.On("GetByID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(cafeIdMatches)).Return(
			testCase.oldCafe, nil)

		cafeMatches := func(c cafeModels.Cafe) bool {
			assert.Equal(t, testCase.outputCafe, c, message)
			return c == testCase.outputCafe
		}

		mockCafeRepo.On("Update",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(cafeMatches)).Return(
			testCase.outputCafe, nil)

		session := sessions.Session{Values: map[interface{}]interface{}{"userID": testCase.staff.StaffID}}
		c := context.WithValue(context.Background(), "session", &session)

		realCafee, err := cafeUsecase.Update(c, testCase.inputCafe)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.outputCafe, realCafee, message)
		}
	}
}
