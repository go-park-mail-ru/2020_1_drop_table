package http

import (
	passMocks "2020_1_drop_table/internal/app/apple_passkit/mocks"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"2020_1_drop_table/internal/pkg/responses"
	"bytes"
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

func TestAddHandler(t *testing.T) {
	type passHttpResponse struct {
		Data   models.UpdateResponse
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		cafeID    string
		publish   string
		inputPass *models.ApplePassDB
		data      models.UpdateResponse
		httpErrs  []responses.HttpError
		err       error
	}

	var cafeID int
	err := faker.FakeData(&cafeID)
	assert.NoError(t, err)

	var inputPass models.ApplePassDB
	err = faker.FakeData(&inputPass)
	assert.NoError(t, err)

	var data models.UpdateResponse
	err = faker.FakeData(&data)
	assert.NoError(t, err)

	url := fmt.Sprintf("/api/v1/cafe/%d/apple_pass", cafeID)

	testCases := []addCafeHandlerTestCase{
		//Test OK
		{
			cafeID:    strconv.Itoa(cafeID),
			publish:   "true",
			inputPass: &inputPass,
			data:      data,
			httpErrs:  nil,
			err:       nil,
		},
		//Test bool parse error
		{
			cafeID:    strconv.Itoa(cafeID),
			publish:   "NOT BOOL",
			inputPass: &inputPass,
			data:      models.UpdateResponse{},
			httpErrs: []responses.HttpError{{
				Code:    400,
				Message: `strconv.ParseBool: parsing "NOT BOOL": invalid syntax`,
			},
			},
			err: nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassUcase := new(passMocks.Usecase)
		handler := applePassKitHandler{passesUsecace: mockPassUcase}

		mockPassUcase.On("UpdatePass",
			mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("models.ApplePassDB")).
			Return(testCase.data, testCase.err)

		var buf bytes.Buffer
		var wr *multipart.Writer
		if testCase.inputPass != nil {
			requestData, err := json.Marshal(&testCase.inputPass)
			assert.NoError(t, err, message)
			buf, wr = createMultipartFormData(t, string(requestData))
		} else {
			buf, wr = createMultipartFormData(t, "")
		}

		req, err := http.NewRequest(echo.PUT, url, &buf)
		assert.NoError(t, err, message)

		req.Header.Set("Content-Type", wr.FormDataContentType())

		req = mux.SetURLVars(req, map[string]string{
			"id": testCase.cafeID,
		})
		req.URL.RawQuery = fmt.Sprintf("publish=%s", testCase.publish)

		respWriter := httptest.NewRecorder()

		handler.UpdatePassHandler(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)

		var responseStruct passHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.NoError(t, err, message)

		errs := responseStruct.Errors
		pass := responseStruct.Data

		assert.Equal(t, testCase.httpErrs, errs, message)
		assert.Equal(t, testCase.data, pass, message)
	}
}

func TestGetPassHandler(t *testing.T) {
	type passHttpResponse struct {
		Data   map[string]string
		Errors []responses.HttpError
	}

	type addCafeHandlerTestCase struct {
		cafeID   string
		publish  string
		data     map[string]string
		httpErrs []responses.HttpError
		err      error
	}

	var cafeID int
	err := faker.FakeData(&cafeID)
	assert.NoError(t, err)

	var data map[string]string
	err = faker.FakeData(&data)
	assert.NoError(t, err)
	url := fmt.Sprintf("/api/v1/cafe/%d/apple_pass", cafeID)

	testCases := []addCafeHandlerTestCase{
		//Test OK
		{
			cafeID:   strconv.Itoa(cafeID),
			publish:  "true",
			data:     data,
			httpErrs: nil,
		},
		{
			cafeID:  "not int",
			publish: "true",
			data:    nil,
			httpErrs: []responses.HttpError{
				{
					Code:    400,
					Message: "bad id: not int",
				},
			},
		},
		{
			cafeID:  strconv.Itoa(cafeID),
			publish: "not bool",
			data:    nil,
			httpErrs: []responses.HttpError{
				{
					Code:    400,
					Message: "strconv.ParseBool: parsing \"not bool\": invalid syntax",
				},
			},
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassUcase := new(passMocks.Usecase)
		handler := applePassKitHandler{passesUsecace: mockPassUcase}

		mockPassUcase.On("GetPass",
			mock.AnythingOfType("*context.valueCtx"),
			mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
			Return(testCase.data, testCase.err)

		req, err := http.NewRequest(echo.GET, url, nil)
		assert.NoError(t, err, message)

		req = mux.SetURLVars(req, map[string]string{
			"id": testCase.cafeID,
		})
		req.URL.RawQuery = fmt.Sprintf("published=%s", testCase.publish)

		respWriter := httptest.NewRecorder()

		handler.GetPassHandler(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)

		var responseStruct passHttpResponse
		err = json.Unmarshal(body, &responseStruct)
		assert.NoError(t, err, message)

		errs := responseStruct.Errors
		pass := responseStruct.Data

		assert.Equal(t, testCase.httpErrs, errs, message)
		assert.Equal(t, testCase.data, pass, message)
	}
}

func TestApplePassKitHandler_GetImageHandler(t *testing.T) {
	type passHttpResponse struct {
		Data   []byte
		Errors []responses.HttpError
	}

	type getImageHandlerTestCase struct {
		imageName         string
		CafeID            string
		loyaltySystemType string
		published         string
		err               error
		httpError         bool
		httpErrors        passHttpResponse
		response          []byte
	}

	var cafeID int
	err := faker.FakeData(&cafeID)
	assert.NoError(t, err)

	var data map[string]string
	err = faker.FakeData(&data)
	assert.NoError(t, err)

	var fakeImage []byte
	err = faker.FakeData(&fakeImage)
	assert.NoError(t, err)

	testCases := []getImageHandlerTestCase{
		//Test OK
		{
			imageName:         "icon@2x.png",
			CafeID:            strconv.Itoa(cafeID),
			loyaltySystemType: "coffee_cup",
			published:         "true",
			err:               nil,
			response:          fakeImage,
		},
		//Test err cafeID not int
		{
			imageName:         "icon@2x.png",
			CafeID:            "not int",
			loyaltySystemType: "coffee_cup",
			published:         "true",
			err:               nil,
			httpError:         true,
			httpErrors: passHttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "bad id: not int",
					},
				},
			},
		},
		//Test err published not bool
		{
			imageName:         "icon@2x.png",
			CafeID:            strconv.Itoa(cafeID),
			loyaltySystemType: "coffee_cup",
			published:         "not bool",
			err:               nil,
			httpError:         true,
			httpErrors: passHttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "strconv.ParseBool: parsing \"not bool\": invalid syntax",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassUcase := new(passMocks.Usecase)
		handler := applePassKitHandler{passesUsecace: mockPassUcase}

		mockPassUcase.On("GetImage",
			mock.AnythingOfType("*context.valueCtx"),
			mock.AnythingOfType("string"), mock.AnythingOfType("int"),
			mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
			Return(testCase.response, testCase.err)

		url := fmt.Sprintf("/api/v1/cafe/%s/apple_pass/%s/%s",
			testCase.CafeID, testCase.loyaltySystemType, testCase.imageName)

		req, err := http.NewRequest(echo.GET, url, nil)
		assert.NoError(t, err, message)

		req = mux.SetURLVars(req, map[string]string{
			"id":                  testCase.CafeID,
			"image_name":          testCase.imageName,
			"loyalty_system_type": testCase.loyaltySystemType,
		})
		req.URL.RawQuery = fmt.Sprintf("published=%s", testCase.published)

		respWriter := httptest.NewRecorder()

		handler.GetImageHandler(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)

		if testCase.httpError {
			var response passHttpResponse
			err = json.Unmarshal(body, &response)
			assert.NoError(t, err, message)
			assert.Equal(t, testCase.httpErrors, response, message)
		} else {
			fileName := fmt.Sprintf("attachment; filename=%s.png", testCase.imageName)
			assert.Equal(t, resp.Header.Get("Content-Disposition"), fileName)
			assert.Equal(t, resp.Header.Get("Content-Type"), "image/png")
			assert.Equal(t, body, testCase.response)
		}
	}
}

func TestApplePassKitHandler_GenerateNewPass(t *testing.T) {
	type passHttpResponse struct {
		Data   []byte
		Errors []responses.HttpError
	}

	type getImageHandlerTestCase struct {
		CafeID            string
		loyaltySystemType string
		published         string
		err               error
		httpError         bool
		httpErrors        passHttpResponse
		response          *bytes.Buffer
	}

	var cafeID int
	err := faker.FakeData(&cafeID)
	assert.NoError(t, err)

	var data map[string]string
	err = faker.FakeData(&data)
	assert.NoError(t, err)

	var fakeImage []byte
	err = faker.FakeData(&fakeImage)
	assert.NoError(t, err)

	newBuffer := new(bytes.Buffer)
	newBuffer.Write(fakeImage)

	testCases := []getImageHandlerTestCase{
		//Test OK
		{
			CafeID:            strconv.Itoa(cafeID),
			loyaltySystemType: "coffee_cup",
			published:         "true",
			err:               nil,
			response:          newBuffer,
		},
		//Test err cafeID not int
		{
			CafeID:            "not int",
			loyaltySystemType: "coffee_cup",
			published:         "true",
			err:               nil,
			httpError:         true,
			httpErrors: passHttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "bad id: not int",
					},
				},
			},
		},
		//Test err published not bool
		{
			CafeID:            strconv.Itoa(cafeID),
			loyaltySystemType: "coffee_cup",
			published:         "not bool",
			err:               nil,
			httpError:         true,
			httpErrors: passHttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "strconv.ParseBool: parsing \"not bool\": invalid syntax",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassUcase := new(passMocks.Usecase)
		handler := applePassKitHandler{passesUsecace: mockPassUcase}

		mockPassUcase.On("GeneratePassObject",
			mock.AnythingOfType("*context.valueCtx"),
			mock.AnythingOfType("int"), mock.AnythingOfType("string"),
			mock.AnythingOfType("bool")).Return(testCase.response, testCase.err)

		url := fmt.Sprintf("/api/v1/cafe/%s/apple_pass/%s/new_customer",
			testCase.CafeID, testCase.loyaltySystemType)

		req, err := http.NewRequest(echo.GET, url, nil)
		assert.NoError(t, err, message)

		req = mux.SetURLVars(req, map[string]string{
			"id":                  testCase.CafeID,
			"loyalty_system_type": testCase.loyaltySystemType,
		})
		req.URL.RawQuery = fmt.Sprintf("published=%s", testCase.published)

		respWriter := httptest.NewRecorder()

		handler.GenerateNewPass(respWriter, req)

		resp := respWriter.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, message)

		if testCase.httpError {
			var response passHttpResponse
			err = json.Unmarshal(body, &response)
			assert.NoError(t, err, message)
			assert.Equal(t, testCase.httpErrors, response, message)
		} else {
			headerValue := "attachment; filename=loyaltyCard.pkpass"
			assert.Equal(t, resp.Header.Get("Content-Disposition"), headerValue)
			assert.Equal(t, resp.Header.Get("Content-Type"), "application/vnd.apple.pkpass")
			assert.Equal(t, body, testCase.response.Bytes())
		}
	}
}
