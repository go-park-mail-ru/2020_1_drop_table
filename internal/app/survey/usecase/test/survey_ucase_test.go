package test

import (
	cafeMock "2020_1_drop_table/internal/app/cafe/mocks"
	models2 "2020_1_drop_table/internal/app/cafe/models"
	globalModels "2020_1_drop_table/internal/app/models"
	staffMock "2020_1_drop_table/internal/app/staff/mocks"
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/app/survey/mocks"
	"2020_1_drop_table/internal/app/survey/usecase"
	"context"
	"github.com/stretchr/testify/mock"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetTemplate(t *testing.T) {

	type CheckStructInput struct {
		Ctx    context.Context
		Survey string
		CafeID int
	}
	type CheckStructOutput struct {
		Err error
	}

	type testCaseStruct struct {
		InputData         CheckStructInput
		OutputData        CheckStructOutput
		RetGetFromContext models.SafeStaff
		RetGetCafeByOwner models2.Cafe
	}

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	c := context.WithValue(context.Background(), "session", &session)

	session2 := sessions.Session{Values: map[interface{}]interface{}{"userID": 1}}
	c2 := context.WithValue(context.Background(), "session", &session2)

	testCases := []testCaseStruct{
		//all ok case
		{
			InputData: CheckStructInput{
				Ctx:    c,
				Survey: "{}",
				CafeID: 229,
			},
			OutputData: CheckStructOutput{
				Err: nil,
			},
			RetGetFromContext: models.SafeStaff{
				StaffID:  228,
				Name:     "",
				Email:    "",
				EditedAt: time.Time{},
				Photo:    "",
				IsOwner:  true,
				CafeId:   0,
				Position: "",
			},
			RetGetCafeByOwner: models2.Cafe{
				CafeID:      1,
				CafeName:    "",
				Address:     "",
				Description: "",
				StaffID:     228,
				OpenTime:    time.Time{},
				CloseTime:   time.Time{},
				Photo:       "",
			},
		},
		{
			InputData: CheckStructInput{
				Ctx:    c2,
				Survey: "Not valid Survey",
				CafeID: -1, //not in
			},
			OutputData: CheckStructOutput{
				Err: globalModels.ErrForbidden,
			},
		},
	}

	timeout := time.Second * 4
	surveyRepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	staffUsecase := staffMock.Usecase{}
	s := usecase.NewSurveyUsecase(&cafeRepo, &surveyRepo, &staffUsecase, timeout)

	for _, testCase := range testCases {
		staffUsecase.On("GetFromSession", mock.AnythingOfType("*context.timerCtx")).Return(testCase.RetGetFromContext, nil)
		cafeRepo.On("GetByID", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.CafeID).Return(testCase.RetGetCafeByOwner, nil)
		surveyRepo.On("SetSurveyTemplate", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.Survey, testCase.InputData.CafeID, testCase.RetGetFromContext.StaffID).Return(nil)
		errRes := s.SetSurveyTemplate(testCase.InputData.Ctx, testCase.InputData.Survey, testCase.InputData.CafeID)
		assert.Equal(t, testCase.OutputData.Err, errRes)
	}
}
