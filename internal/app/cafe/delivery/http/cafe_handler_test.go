package http_test

import (
	cafeHandlers "2020_1_drop_table/internal/app/cafe/delivery/http"
	cafeMocks "2020_1_drop_table/internal/app/cafe/mocks"
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/pkg/responses"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/gorilla/mux"
	"github.com/labstack/echo"
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

const url = "/api/v1/cafe"

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

func TestAddCafeHandler(t *testing.T) {
	type cafeHttpResponse struct {
		Data   cafeModels.Cafe
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		inputCafe  *cafeModels.Cafe
		outputCafe cafeModels.Cafe
		httpErrs   []responses.HttpError
	}

	mockCafeUcase := new(cafeMocks.Usecase)
	handler := cafeHandlers.CafeHandler{CUsecase: mockCafeUcase}

	var inputCafe cafeModels.Cafe
	err := faker.FakeData(&inputCafe)
	assert.NoError(t, err)

	inputCafe.CloseTime = inputCafe.CloseTime.UTC()
	inputCafe.OpenTime = inputCafe.OpenTime.UTC()

	outputCafe := inputCafe

	inputCafe.CafeID = 0

	testCases := []addCafeHandlerTestCase{
		//Test OK
		{
			inputCafe:  &inputCafe,
			outputCafe: outputCafe,
			httpErrs:   nil,
		},
		//Test empty JsonData
		{
			inputCafe:  nil,
			outputCafe: cafeModels.Cafe{},
			httpErrs: []responses.HttpError{{
				Code:    400,
				Message: globalModels.ErrEmptyJSON.Error(),
			},
			},
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		cafeMatches := func(c cafeModels.Cafe) bool {
			assert.Equal(t, *testCase.inputCafe, c, message)
			return c == *testCase.inputCafe
		}

		mockCafeUcase.On("Add",
			mock.AnythingOfType("*context.emptyCtx"), mock.MatchedBy(cafeMatches)).
			Return(testCase.outputCafe, nil)

		var buf bytes.Buffer
		var wr *multipart.Writer
		if testCase.inputCafe != nil {
			requestData, err := json.Marshal(&testCase.inputCafe)
			assert.NoError(t, err, message)
			buf, wr = createMultipartFormData(t, string(requestData))
		} else {
			buf, wr = createMultipartFormData(t, "")
		}

		req, err := http.NewRequest(echo.POST, url, &buf)
		assert.NoError(t, err, message)
		req.Header.Set("Content-Type", wr.FormDataContentType())

		respWriter := httptest.NewRecorder()

		handler.AddCafeHandler(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)

		var responseStruct cafeHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.NoError(t, err, message)

		errs := responseStruct.Errors
		cafe := testCase.outputCafe

		assert.Equal(t, testCase.httpErrs, errs, message)
		assert.Equal(t, testCase.outputCafe, cafe, message)
	}
}

func TestEditCafeHandler(t *testing.T) {
	type cafeHttpResponse struct {
		Data   cafeModels.Cafe
		Errors []responses.HttpError
	}

	type UpdateTestCase struct {
		cafeID     string
		inputCafe  *cafeModels.Cafe
		outputCafe cafeModels.Cafe
		httpErrs   []responses.HttpError
		sqlErr     error
	}

	var inputCafe cafeModels.Cafe
	err := faker.FakeData(&inputCafe)
	assert.NoError(t, err)

	inputCafe.CloseTime = inputCafe.CloseTime.UTC()
	inputCafe.OpenTime = inputCafe.OpenTime.UTC()

	outputCafe := inputCafe

	inputCafe.CafeID = 0

	testCases := []UpdateTestCase{
		//Test OK
		{
			cafeID:     strconv.Itoa(outputCafe.CafeID),
			inputCafe:  &inputCafe,
			outputCafe: outputCafe,
			httpErrs:   nil,
			sqlErr:     nil,
		},
		//Test empty JsonData
		{
			inputCafe:  nil,
			outputCafe: cafeModels.Cafe{},
			httpErrs: []responses.HttpError{{
				Code:    400,
				Message: globalModels.ErrEmptyJSON.Error(),
			},
			},
		},
		//Test not found
		{
			cafeID:     strconv.Itoa(outputCafe.CafeID + 1),
			inputCafe:  &cafeModels.Cafe{},
			outputCafe: cafeModels.Cafe{CafeID: outputCafe.CafeID + 1},
			sqlErr:     sql.ErrNoRows,
			httpErrs: []responses.HttpError{
				{
					Code:    400,
					Message: sql.ErrNoRows.Error(),
				},
			},
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockCafeUcase := new(cafeMocks.Usecase)
		handler := cafeHandlers.CafeHandler{CUsecase: mockCafeUcase}

		cafeEqual := func(c cafeModels.Cafe) bool {
			assert.Equal(t, c, testCase.outputCafe, message)
			return c == testCase.outputCafe
		}

		mockCafeUcase.On("Update",
			mock.AnythingOfType("*context.valueCtx"), mock.MatchedBy(cafeEqual)).
			Return(testCase.outputCafe, testCase.sqlErr)

		var buf bytes.Buffer
		var wr *multipart.Writer
		if testCase.inputCafe != nil {
			requestData, err := json.Marshal(&testCase.inputCafe)
			assert.NoError(t, err, message)
			buf, wr = createMultipartFormData(t, string(requestData))
		} else {
			buf, wr = createMultipartFormData(t, "")
		}

		req, err := http.NewRequest(echo.POST, url, &buf)
		assert.NoError(t, err, message)

		req.Header.Set("Content-Type", wr.FormDataContentType())

		req = mux.SetURLVars(req, map[string]string{
			"id": testCase.cafeID,
		})

		respWriter := httptest.NewRecorder()

		handler.EditCafeHandler(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)

		var responseStruct cafeHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.NoError(t, err, message)

		errs := responseStruct.Errors
		cafe := testCase.outputCafe

		assert.Equal(t, testCase.httpErrs, errs, message)
		assert.Equal(t, testCase.outputCafe, cafe, message)
	}

}

func TestGetByIDHandler(t *testing.T) {
	type cafeHttpResponse struct {
		Data   cafeModels.Cafe
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		cafeID     string
		outputCafe cafeModels.Cafe
		httpErrs   []responses.HttpError
		sqlErr     error
	}

	cafeID := 1

	var outputCafe cafeModels.Cafe
	err := faker.FakeData(&outputCafe)
	assert.NoError(t, err)

	outputCafe.CloseTime = outputCafe.CloseTime.UTC()
	outputCafe.OpenTime = outputCafe.OpenTime.UTC()
	outputCafe.CafeID = cafeID

	testCases := []addCafeHandlerTestCase{
		//Test OK
		{
			cafeID:     strconv.Itoa(cafeID),
			outputCafe: outputCafe,
			httpErrs:   nil,
		},
		//Test int parsing error
		{
			cafeID:     "NOT INT",
			outputCafe: outputCafe,
			httpErrs: []responses.HttpError{
				{
					Code:    400,
					Message: "bad id: NOT INT",
				},
			},
		},
		//Test not found
		{
			cafeID:     strconv.Itoa(cafeID + 1),
			outputCafe: cafeModels.Cafe{},
			sqlErr:     sql.ErrNoRows,
			httpErrs: []responses.HttpError{
				{
					Code:    400,
					Message: sql.ErrNoRows.Error(),
				},
			},
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockCafeUcase := new(cafeMocks.Usecase)
		handler := cafeHandlers.CafeHandler{CUsecase: mockCafeUcase}

		idMatches := func(id int) bool {
			testCaseID, err := strconv.Atoi(testCase.cafeID)
			if err != nil {
				return false
			}
			assert.Equal(t, testCaseID, id, message)
			return id == testCaseID
		}

		mockCafeUcase.On("GetByID",
			mock.AnythingOfType("*context.valueCtx"), mock.MatchedBy(idMatches)).
			Return(testCase.outputCafe, testCase.sqlErr)

		req, err := http.NewRequest(echo.GET, url, nil)
		assert.NoError(t, err, message)
		req = mux.SetURLVars(req, map[string]string{
			"id": testCase.cafeID,
		})

		respWriter := httptest.NewRecorder()

		handler.GetByIDHandler(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)

		var responseStruct cafeHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.NoError(t, err, message)

		errs := responseStruct.Errors
		cafe := testCase.outputCafe

		assert.Equal(t, testCase.httpErrs, errs, message)
		assert.Equal(t, testCase.outputCafe, cafe, message)
	}
}

func TestGetByOwnerIDHandler(t *testing.T) {
	type cafeHttpResponse struct {
		Data   []cafeModels.Cafe
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		outputCafe []cafeModels.Cafe
		httpErrs   []responses.HttpError
		sqlErr     error
	}

	outputCafes := make([]cafeModels.Cafe, 7)
	err := faker.FakeData(&outputCafes)
	assert.NoError(t, err)

	for i := range outputCafes {
		outputCafes[i].CloseTime = outputCafes[i].CloseTime.UTC()
		outputCafes[i].OpenTime = outputCafes[i].OpenTime.UTC()
	}

	testCases := []addCafeHandlerTestCase{
		//Test OK
		{
			outputCafe: outputCafes,
			httpErrs:   nil,
			sqlErr:     nil,
		},
		//Test not Found
		{
			outputCafe: nil,
			httpErrs: []responses.HttpError{
				{
					Code:    400,
					Message: sql.ErrNoRows.Error(),
				},
			},
			sqlErr: sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockCafeUcase := new(cafeMocks.Usecase)
		handler := cafeHandlers.CafeHandler{CUsecase: mockCafeUcase}

		mockCafeUcase.On("GetByOwnerID",
			mock.AnythingOfType("*context.emptyCtx")).Return(testCase.outputCafe, testCase.sqlErr)

		req, err := http.NewRequest(echo.GET, url, nil)
		assert.NoError(t, err, message)

		respWriter := httptest.NewRecorder()

		handler.GetByOwnerIDHandler(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)
		var responseStruct cafeHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.NoError(t, err, message)

		errs := responseStruct.Errors
		cafes := testCase.outputCafe

		assert.Equal(t, testCase.httpErrs, errs, message)
		assert.Equal(t, testCase.outputCafe, cafes, message)
	}
}
