package usecase

import (
	cafeMock "2020_1_drop_table/internal/app/cafe/mocks"
	models2 "2020_1_drop_table/internal/app/cafe/models"
	globalModels "2020_1_drop_table/internal/app/models"
	staffMock "2020_1_drop_table/internal/microservices/staff/mocks"
	"2020_1_drop_table/internal/microservices/staff/models"
	"2020_1_drop_table/internal/microservices/survey/mocks"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		InputData            CheckStructInput
		OutputData           CheckStructOutput
		RetGetFromContext    models.SafeStaff
		RetGetCafeByOwner    models2.Cafe
		SetSurveyTemplateErr error
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
		//update case
		{
			InputData: CheckStructInput{
				Ctx:    c,
				Survey: `{"test":"value"}`,
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
			SetSurveyTemplateErr: errors.New(`pq: duplicate key value violates unique constraint "surveytemplate_cafeid_key"`),
		},
	}

	timeout := time.Second * 4
	surveyRepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	staffUsecase := staffMock.Usecase{}
	s := NewSurveyUsecase(&surveyRepo, timeout)

	for _, testCase := range testCases {
		staffUsecase.On("GetFromSession", mock.AnythingOfType("*context.timerCtx")).Return(testCase.RetGetFromContext, nil)
		cafeRepo.On("GetByID", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.CafeID).Return(testCase.RetGetCafeByOwner, nil)
		surveyRepo.On("SetSurveyTemplate", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.Survey, testCase.InputData.CafeID, testCase.RetGetFromContext.StaffID).Return(testCase.SetSurveyTemplateErr)
		surveyRepo.On("UpdateSurveyTemplate", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.Survey, testCase.InputData.CafeID).Return(nil)
		errRes := s.SetSurveyTemplate(testCase.InputData.Ctx, testCase.InputData.Survey, testCase.InputData.CafeID)
		assert.Equal(t, testCase.OutputData.Err, errRes)
	}
}

func TestGetSurvey(t *testing.T) {

	type CheckStructInput struct {
		Ctx    context.Context
		CafeID int
	}
	type CheckStructOutput struct {
		Err    error
		Survey string
	}

	type testCaseStruct struct {
		InputData    CheckStructInput
		OutputData   CheckStructOutput
		GetSurveyErr error
	}

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	c := context.WithValue(context.Background(), "session", &session)

	testCases := []testCaseStruct{

		//all ok case
		{
			InputData: CheckStructInput{
				Ctx:    c,
				CafeID: 229,
			},
			OutputData: CheckStructOutput{
				Err:    nil,
				Survey: `{"test":"value"}`,
			},
		},
		//not found case
		{
			InputData: CheckStructInput{
				Ctx:    c,
				CafeID: 123, //not in
			},
			OutputData: CheckStructOutput{
				Err:    nil,
				Survey: "",
			},
			GetSurveyErr: sql.ErrNoRows,
		},
	}

	timeout := time.Second * 4
	surveyRepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	staffUsecase := staffMock.Usecase{}
	s := NewSurveyUsecase(&cafeRepo, &surveyRepo, &staffUsecase, timeout)

	for _, testCase := range testCases {
		surveyRepo.On("GetSurveyTemplate", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.CafeID).Return(testCase.OutputData.Survey, testCase.GetSurveyErr)
		survey, errRes := s.GetSurveyTemplate(testCase.InputData.Ctx, testCase.InputData.CafeID)
		assert.Equal(t, testCase.OutputData.Err, errRes)
		assert.Equal(t, testCase.OutputData.Survey, survey)
	}
}

func TestSubmitSurvey(t *testing.T) {
	type CheckStructInput struct {
		Ctx          context.Context
		Survey       string
		CustomerUUID string
	}
	type CheckStructOutput struct {
		Err error
	}

	type testCaseStruct struct {
		InputData       CheckStructInput
		OutputData      CheckStructOutput
		SubmitSurveyErr error
	}

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	c := context.WithValue(context.Background(), "session", &session)

	testCases := []testCaseStruct{

		//all ok case
		{
			InputData: CheckStructInput{
				Ctx:          c,
				Survey:       "{}",
				CustomerUUID: "good uuid",
			},
			OutputData: CheckStructOutput{
				Err: nil,
			},
			SubmitSurveyErr: nil,
		},
		{
			InputData: CheckStructInput{
				Ctx:          c,
				Survey:       "Not valid Survey",
				CustomerUUID: "Not valid UUID", //not in
			},
			OutputData: CheckStructOutput{
				Err: globalModels.ErrBadUuid,
			},
			SubmitSurveyErr: errors.New(fmt.Sprintf(`pq: invalid input syntax for type uuid: "%s"`, "Not valid UUID")),
		},
	}

	timeout := time.Second * 4
	surveyRepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	staffUsecase := staffMock.Usecase{}
	s := NewSurveyUsecase(&cafeRepo, &surveyRepo, &staffUsecase, timeout)

	for _, testCase := range testCases {
		surveyRepo.On("SubmitSurvey", mock.AnythingOfType("*context.timerCtx"), testCase.InputData.Survey, testCase.InputData.CustomerUUID).Return(testCase.SubmitSurveyErr)
		errRes := s.SubmitSurvey(testCase.InputData.Ctx, testCase.InputData.Survey, testCase.InputData.CustomerUUID)
		assert.Equal(t, testCase.OutputData.Err, errRes)
	}
}
