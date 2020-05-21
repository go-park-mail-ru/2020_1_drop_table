package usecase

import (
	"2020_1_drop_table/configs"
	cafeClientGRPCMock "2020_1_drop_table/internal/app/cafe/delivery/grpc/client/mocks"
	models3 "2020_1_drop_table/internal/app/cafe/models"
	"2020_1_drop_table/internal/app/statistics/mocks"
	"2020_1_drop_table/internal/app/statistics/models"
	staffClientGRPCMock "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client/mocks"
	models2 "2020_1_drop_table/internal/microservices/staff/models"
	"context"
	"database/sql"
	"errors"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGetSurvey(t *testing.T) {

	type CheckStructInput struct {
		Ctx     context.Context
		StaffID int
		Limit   int
		Since   int
	}
	type CheckStructOutput struct {
		Stat []models.StatisticsStruct
		Err  error
	}

	type testCaseStruct struct {
		InputData           CheckStructInput
		OutputData          CheckStructOutput
		StaffFromSession    models2.SafeStaff
		StaffFromSessionErr error
	}

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	c := context.WithValue(context.Background(), configs.SessionStaffID, &session)

	firstStatData := []models.StatisticsStruct{{
		JsonData:   "valid",
		Time:       time.Time{},
		ClientUUID: "asdasd",
		StaffId:    229,
		CafeId:     2,
	}}

	testCases := []testCaseStruct{

		//all ok case
		{
			InputData: CheckStructInput{
				Ctx:     c,
				StaffID: 229,
				Limit:   5,
				Since:   0,
			},
			OutputData: CheckStructOutput{
				Err:  nil,
				Stat: firstStatData,
			},
			StaffFromSession: models2.SafeStaff{
				StaffID:  123,
				Name:     "",
				Email:    "",
				EditedAt: time.Time{},
				Photo:    "",
				IsOwner:  true,
				CafeId:   2,
				Position: "",
			},
		},
	}

	timeout := time.Second * 4
	statRepo := new(mocks.Repository)
	cafeRepo := new(cafeClientGRPCMock.CafeGRPCClientInterface)
	staffUsecase := new(staffClientGRPCMock.StaffClientInterface)
	s := NewStatisticsUsecase(statRepo, staffUsecase, cafeRepo, timeout)

	for _, testCase := range testCases {
		staffUsecase.On("GetFromSession", mock.AnythingOfType("*context.timerCtx")).Return(testCase.StaffFromSession, testCase.StaffFromSessionErr)
		statRepo.On("GetWorkerDataFromRepo", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.StaffID, testCase.InputData.Limit, testCase.InputData.Since).Return(testCase.OutputData.Stat, testCase.OutputData.Err)
		stat, err := s.GetWorkerData(testCase.InputData.Ctx, testCase.InputData.StaffID, testCase.InputData.Limit, testCase.InputData.Since)
		assert.Equal(t, testCase.OutputData.Err, err)
		assert.Equal(t, testCase.OutputData.Stat, stat)
	}
}

func TestGetDataFroGraphs(t *testing.T) {

	type CheckStructInput struct {
		Ctx   context.Context
		Typ   string
		Since string
		To    string
	}
	type CheckStructOutput struct {
		Stat map[string]map[string][]models.TempStruct
		Err  error
	}

	type testCaseStruct struct {
		InputData                 CheckStructInput
		OutputData                CheckStructOutput
		StaffFromSession          models2.SafeStaff
		StaffFromSessionErr       error
		GetCafesCafes             []models3.Cafe
		GetCafesErr               error
		GetGraphsDataFromRepoData []models.StatisticsGraphRawStruct
		GetGraphsDataFromRepoErr  error
	}

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	c := context.WithValue(context.Background(), configs.SessionStaffID, &session)
	tNow := time.Now()
	emptStat:=map[string]map[string][]models.TempStruct(nil)
	res := map[string]map[string][]models.TempStruct{
		"1": map[string][]models.TempStruct{
			"1": []models.TempStruct{models.TempStruct{NumOfUsage: 10, Date: tNow},
				models.TempStruct{NumOfUsage: 7, Date: tNow}},
			"2": []models.TempStruct{models.TempStruct{NumOfUsage: 8, Date: tNow}}},
		"2": map[string][]models.TempStruct{"3": []models.TempStruct{models.TempStruct{NumOfUsage: 8, Date: tNow}}}}
	testCases := []testCaseStruct{

		//all ok case
		{
			InputData: CheckStructInput{
				Ctx:   c,
				Typ:   "day",
				To:    time.Now().String(),
				Since: time.Now().String(),
			},
			OutputData: CheckStructOutput{
				Err:  nil,
				Stat: res,
			},
			StaffFromSession: models2.SafeStaff{
				StaffID:  123,
				Name:     "",
				Email:    "",
				EditedAt: time.Time{},
				Photo:    "",
				IsOwner:  true,
				CafeId:   0,
				Position: "",
			},
			GetCafesCafes: []models3.Cafe{
				{
					CafeID:      1,
					CafeName:    "",
					Address:     "",
					Description: "",
					StaffID:     123,
					OpenTime:    time.Time{},
					CloseTime:   time.Time{},
					Photo:       "",
					Location:    "",
				},
				{
					CafeID:      2,
					CafeName:    "",
					Address:     "",
					Description: "",
					StaffID:     123,
					OpenTime:    time.Time{},
					CloseTime:   time.Time{},
					Photo:       "",
					Location:    "",
				},
			},
			GetCafesErr: nil,
			GetGraphsDataFromRepoData: []models.StatisticsGraphRawStruct{
				{
					Count:   10,
					Date:    tNow,
					CafeId:  1,
					StaffId: 1,
				},
				{
					Count:   7,
					Date:    tNow,
					CafeId:  1,
					StaffId: 1,
				},
				{
					Count:   8,
					Date:    tNow,
					CafeId:  1,
					StaffId: 2,
				},
				{
					Count:   8,
					Date:    tNow,
					CafeId:  2,
					StaffId: 3,
				},
			},
		},
		//not ok case
		{
			InputData: CheckStructInput{
				Ctx:   c,
				Typ:   "day",
				To:    time.Now().String(),
				Since: time.Now().String(),
			},
			OutputData: CheckStructOutput{
				Err:  errors.New("cant found statistics data with this input"),
				Stat: emptStat,
			},
			StaffFromSession: models2.SafeStaff{
				StaffID:  228,
				Name:     "",
				Email:    "",
				EditedAt: time.Time{},
				Photo:    "",
				IsOwner:  true,
				CafeId:   0,
				Position: "",
			},
			GetCafesCafes: []models3.Cafe{
				{
					CafeID:      1,
					CafeName:    "",
					Address:     "",
					Description: "",
					StaffID:     123,
					OpenTime:    time.Time{},
					CloseTime:   time.Time{},
					Photo:       "",
					Location:    "",
				},
				{
					CafeID:      2,
					CafeName:    "",
					Address:     "",
					Description: "",
					StaffID:     123,
					OpenTime:    time.Time{},
					CloseTime:   time.Time{},
					Photo:       "",
					Location:    "",
				},
			},
			GetCafesErr:   nil,
			GetGraphsDataFromRepoErr: sql.ErrNoRows,

		},
	}

	timeout := time.Second * 4
	statRepo := new(mocks.Repository)
	cafeRepo := new(cafeClientGRPCMock.CafeGRPCClientInterface)
	staffUsecase := new(staffClientGRPCMock.StaffClientInterface)
	s := NewStatisticsUsecase(statRepo, staffUsecase, cafeRepo, timeout)
	for _, testCase := range testCases {
		staffUsecase.On("GetFromSession", mock.AnythingOfType("*context.timerCtx")).Return(testCase.StaffFromSession, testCase.StaffFromSessionErr)
		cafeRepo.On("GetByOwnerId", mock.AnythingOfType("*context.timerCtx"), testCase.StaffFromSession.StaffID).Return(testCase.GetCafesCafes, testCase.GetCafesErr)
		statRepo.On("GetGraphsDataFromRepo", mock.AnythingOfType("*context.timerCtx"), testCase.GetCafesCafes, testCase.InputData.Typ, testCase.InputData.Since, testCase.InputData.To).Return(testCase.GetGraphsDataFromRepoData, testCase.GetGraphsDataFromRepoErr)
		stat, err := s.GetDataForGraphs(testCase.InputData.Ctx, testCase.InputData.Typ, testCase.InputData.Since, testCase.InputData.To)
		assert.Equal(t, testCase.OutputData.Err, err)
		assert.Equal(t, testCase.OutputData.Stat, stat)
	}
}
