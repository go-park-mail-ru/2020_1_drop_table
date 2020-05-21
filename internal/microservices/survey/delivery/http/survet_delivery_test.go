package http

import (
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/microservices/survey/mocks"
	"2020_1_drop_table/internal/pkg/responses"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func createMultipartFormData(t *testing.T, data string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)

	var fw io.Writer
	dataReader := strings.NewReader(data)
	if fw, err = w.CreateFormField("jsonData"); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, dataReader); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}

	err = w.Close()
	if err != nil {
		t.Error(err)
	}

	return b, w
}

func TestSetSurveyTemplate(t *testing.T) {
	url := `/api/v1/survey/set_survey_template/{id:[0-9]+}`
	type SurveyHttpResponse struct {
		Data   error
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		inputSurvey string
		CafeID      int
		outputError error
		httpErrs    []responses.HttpError
	}

	mockSurveyUcase := new(mocks.Usecase)
	handler := SurveyHandler{SurveyUC: mockSurveyUcase}

	testCases := []addCafeHandlerTestCase{
		//Test not Valid
		{
			inputSurvey: `ahsjgdhjas""{}`,
			CafeID:      228,
			outputError: globalModels.ErrBadJSON,
			httpErrs: []responses.HttpError{{
				Code:    400,
				Message: globalModels.ErrBadJSON.Error(),
			},
			},
		},
		//Test OK
		{
			inputSurvey: `{}`,
			CafeID:      228,
			outputError: nil,
			httpErrs:    nil,
		},
	}

	for _, testCase := range testCases {
		mockSurveyUcase.On("SetSurveyTemplate",
			mock.AnythingOfType("*context.valueCtx"), testCase.inputSurvey, testCase.CafeID).
			Return(testCase.outputError)

		buf, wr := createMultipartFormData(t, testCase.inputSurvey)
		req, err := http.NewRequest("POST", url, &buf)
		req = mux.SetURLVars(req, map[string]string{
			"id": "228",
		})
		assert.Nil(t, err)
		req.Header.Set("Content-Type", wr.FormDataContentType())
		respWriter := httptest.NewRecorder()
		handler.SetSurveyTemplate(respWriter, req)
		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		var responseStruct SurveyHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.Nil(t, err)
		assert.Equal(t, responseStruct.Errors, testCase.httpErrs)
	}
}

func TestGetSurveyTemplate(t *testing.T) {
	url := `/api/v1/survey/get_survey_template/{id:[0-9]+}`
	type SurveyHttpResponse struct {
		Data   string
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		CafeID       int
		outputError  error
		outputSurvey string
		httpErrs     []responses.HttpError
	}

	mockSurveyUcase := new(mocks.Usecase)
	handler := SurveyHandler{SurveyUC: mockSurveyUcase}

	testCases := []addCafeHandlerTestCase{
		//Test OK
		{
			CafeID:       229,
			outputSurvey: `{}`,
			outputError:  nil,
			httpErrs:     nil,
		},
		//Test not Valid
		{
			CafeID:       228,
			outputSurvey: ``,
			outputError:  globalModels.ErrForbidden,
			httpErrs: []responses.HttpError{{
				Code:    400,
				Message: globalModels.ErrForbidden.Error(),
			},
			},
		},
	}

	for _, testCase := range testCases {
		mockSurveyUcase.On("GetSurveyTemplate",
			mock.AnythingOfType("*context.valueCtx"), testCase.CafeID).
			Return(testCase.outputSurvey, testCase.outputError)

		req, err := http.NewRequest("GET", url, nil)

		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(testCase.CafeID),
		})
		assert.Nil(t, err)
		respWriter := httptest.NewRecorder()
		handler.GetSurveyTemplate(respWriter, req)
		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		var responseStruct SurveyHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.Nil(t, err)
		assert.Equal(t, responseStruct.Data, testCase.outputSurvey)
		assert.Equal(t, responseStruct.Errors, testCase.httpErrs)
	}
}

func TestSubmitSurvey(t *testing.T) {
	url := `/api/v1/survey/submit_survey/{customerid}`
	type SurveyHttpResponse struct {
		Data   error
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		CustomerID  string
		outputError error
		inputSurvey string
		httpErrs    []responses.HttpError
	}

	mockSurveyUcase := new(mocks.Usecase)
	handler := SurveyHandler{SurveyUC: mockSurveyUcase}

	testCases := []addCafeHandlerTestCase{
		//Test not OK
		{
			inputSurvey: `{}`,
			CustomerID:  "not valid",
			outputError: globalModels.ErrForbidden,
			httpErrs: []responses.HttpError{{
				Code:    400,
				Message: globalModels.ErrForbidden.Error(),
			},
			},
		},
		//Test OK
		{
			inputSurvey: `{}`,
			CustomerID:  "asd",
			outputError: nil,
			httpErrs:    []responses.HttpError(nil),
		},
	}

	for _, testCase := range testCases {
		mockSurveyUcase.On("SubmitSurvey",
			mock.AnythingOfType("*context.valueCtx"), testCase.inputSurvey, testCase.CustomerID).
			Return(testCase.outputError)
		buf, wr := createMultipartFormData(t, testCase.inputSurvey)
		req, err := http.NewRequest("POST", url, &buf)
		req = mux.SetURLVars(req, map[string]string{
			"customerid": testCase.CustomerID,
		})
		req.Header.Set("Content-Type", wr.FormDataContentType())
		assert.Nil(t, err)
		respWriter := httptest.NewRecorder()
		handler.SubmitSurvey(respWriter, req)
		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		var responseStruct SurveyHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.Nil(t, err)
		assert.Equal(t, responseStruct.Errors, testCase.httpErrs)
	}

}
