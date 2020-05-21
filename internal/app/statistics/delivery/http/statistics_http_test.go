package http

import (
	"2020_1_drop_table/internal/app/statistics/mocks"
	"2020_1_drop_table/internal/app/statistics/models"
	"2020_1_drop_table/internal/pkg/responses"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetWorkerData(t *testing.T) {
	url := `/api/v1/statistics/get_worker_data`
	type SurveyHttpResponse struct {
		Data   []models.StatisticsStruct
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		inputData    models.GetWorkerDataStruct
		ouputData    []models.StatisticsStruct
		httpErrs     []responses.HttpError
		usecaseError error
	}

	mockStatUcase := new(mocks.Usecase)
	handler := StatisticsHandler{SUsecase: mockStatUcase}
	inputData := models.GetWorkerDataStruct{
		StaffID: 123,
		Since:   0,
		Limit:   0,
	}

	testCases := []addCafeHandlerTestCase{

		//Test OK
		{
			inputData: inputData,
			httpErrs:  nil,
			ouputData: []models.StatisticsStruct{
				{
					JsonData:   "{}",
					Time:       time.Time{},
					ClientUUID: "asd",
					StaffId:    123,
					CafeId:     2,
				},
			},
			usecaseError: nil,
		},
		//Test bad user OK
		{
			inputData: models.GetWorkerDataStruct{
				StaffID: 228,
				Since:   0,
				Limit:   0,
			},
			httpErrs: []responses.HttpError{responses.HttpError{
				Code:    400,
				Message: "sql: no rows in result set",
			}},
			ouputData:    []models.StatisticsStruct{},
			usecaseError: sql.ErrNoRows,
		},
	}

	for _, testCase := range testCases {
		mockStatUcase.On("GetWorkerData",
			mock.AnythingOfType("*context.emptyCtx"), testCase.inputData.StaffID, testCase.inputData.Limit, testCase.inputData.Since).
			Return(testCase.ouputData, testCase.usecaseError)

		buf, _ := json.Marshal(testCase.inputData)
		req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")

		assert.Nil(t, err)
		respWriter := httptest.NewRecorder()
		handler.GetWorkerData(respWriter, req)
		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		var responseStruct SurveyHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.Nil(t, err)
		assert.Equal(t, responseStruct.Errors, testCase.httpErrs)
	}
}

func TestGetDataForGraphs(t *testing.T) {
	url := `/api/v1/statistics/get_graphs_data?type=day&since=-1&to=-1`
	type SurveyHttpResponse struct {
		Data   map[string]map[string][]models.TempStruct
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		typ          string
		since        string
		to           string
		ouputData    map[string]map[string][]models.TempStruct
		httpErrs     []responses.HttpError
		usecaseError error
	}

	mockStatUcase := new(mocks.Usecase)
	handler := StatisticsHandler{SUsecase: mockStatUcase}

	testCases := []addCafeHandlerTestCase{

		//Test not OK
		{
			typ:          "day",
			since:        "-1",
			to:           "-1",
			httpErrs:     []responses.HttpError{responses.HttpError{Code: 400, Message: "sql: no rows in result set"}},
			ouputData:    map[string]map[string][]models.TempStruct{},
			usecaseError: sql.ErrNoRows,
		},
	}

	for _, testCase := range testCases {
		mockStatUcase.On("GetDataForGraphs",
			mock.AnythingOfType("*context.valueCtx"), testCase.typ, testCase.since, testCase.to).
			Return(testCase.ouputData, testCase.usecaseError)

		req, err := http.NewRequest("GET", url, nil)
		req = mux.SetURLVars(req, map[string]string{
			"type":  testCase.typ,
			"since": testCase.since,
			"to":    testCase.to,
		})

		assert.Nil(t, err)
		respWriter := httptest.NewRecorder()
		handler.GetGraphsData(respWriter, req)
		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		var responseStruct SurveyHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.Nil(t, err)
		assert.Equal(t, responseStruct.Errors, testCase.httpErrs)
	}
}
