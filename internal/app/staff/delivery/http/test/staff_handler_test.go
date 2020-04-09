package test

import (
	http2 "2020_1_drop_table/internal/app/staff/delivery/http"
	"2020_1_drop_table/internal/app/staff/mocks"
	"2020_1_drop_table/internal/app/staff/models"
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
	"strings"
	"testing"
	"time"
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

func TestGetById(t *testing.T) {

	returnStaff := models.SafeStaff{
		StaffID:  228,
		Name:     "",
		Email:    "",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  false,
		CafeId:   0,
		Position: "",
	}
	const url = "/api/v1/staff/"

	mockstaffUcase := new(mocks.Usecase)
	mockstaffUcase.On("GetByID", mock.AnythingOfType("*context.valueCtx"), 228).Return(returnStaff, nil)
	handler := http2.StaffHandler{SUsecase: mockstaffUcase}
	buf, wr := createMultipartFormData(t, "")
	req, err := http.NewRequest("GET", url, &buf)
	req = mux.SetURLVars(req, map[string]string{
		"id": "228",
	})
	assert.Nil(t, err)
	req.Header.Set("Content-Type", wr.FormDataContentType())
	respWriter := httptest.NewRecorder()
	handler.GetStaffByIdHandler(respWriter, req)
	resp := respWriter.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var responseStruct staffHttpResponse
	err = json.Unmarshal(body, &responseStruct)
	assert.NoError(t, err)
	assert.Equal(t, responseStruct.Data, returnStaff)

}

type staffHttpResponse struct {
	Data   models.SafeStaff
	Errors []responses.HttpError
}

func TestGetCurrStaff(t *testing.T) {

	const url = "/api/v1/staff"
	mockstaffUcase := new(mocks.Usecase)
	returnStaff := models.SafeStaff{
		StaffID:  228,
		Name:     "",
		Email:    "",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  false,
		CafeId:   0,
		Position: "",
	}
	mockstaffUcase.On("GetFromSession", mock.AnythingOfType("*context.emptyCtx")).Return(returnStaff, nil)
	handler := http2.StaffHandler{SUsecase: mockstaffUcase}
	buf, wr := createMultipartFormData(t, "")
	req, err := http.NewRequest("POST", url, &buf)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", wr.FormDataContentType())
	respWriter := httptest.NewRecorder()

	handler.GetCurrentStaffHandler(respWriter, req)
	resp := respWriter.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var responseStruct staffHttpResponse
	err = json.Unmarshal(body, &responseStruct)
	assert.NoError(t, err)
	assert.Equal(t, responseStruct.Data, returnStaff)

}

func TestGenerateQrHandler(t *testing.T) {
	type QrHttpResponse struct {
		Data   string
		Errors []responses.HttpError
	}
	const url = "/api/v1/staff/generateQr"

	mockstaffUcase := new(mocks.Usecase)
	mockstaffUcase.On("GetQrForStaff", mock.AnythingOfType("*context.valueCtx"), 228).Return("path", nil)
	handler := http2.StaffHandler{SUsecase: mockstaffUcase}
	buf, wr := createMultipartFormData(t, "")
	req, err := http.NewRequest("GET", url, &buf)
	req = mux.SetURLVars(req, map[string]string{
		"id": "228",
	})
	assert.Nil(t, err)
	req.Header.Set("Content-Type", wr.FormDataContentType())
	respWriter := httptest.NewRecorder()
	handler.GenerateQrHandler(respWriter, req)
	resp := respWriter.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var responseStruct QrHttpResponse
	err = json.Unmarshal(body, &responseStruct)
	assert.NoError(t, err)
	assert.Equal(t, responseStruct.Data, "path")

}

func TestEditStaff(t *testing.T) {

	returnStaff := models.SafeStaff{
		StaffID:  228,
		Name:     "",
		Email:    "",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  false,
		CafeId:   0,
		Position: "",
	}
	const url = "/api/v1/staff/"

	mockstaffUcase := new(mocks.Usecase)
	mockstaffUcase.On("Update", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("models.SafeStaff")).Return(returnStaff, nil)
	handler := http2.StaffHandler{SUsecase: mockstaffUcase}
	str, _ := json.Marshal(returnStaff)
	buf, wr := createMultipartFormData(t, string(str))
	req, err := http.NewRequest("GET", url, &buf)
	req = mux.SetURLVars(req, map[string]string{
		"id": "228",
	})
	assert.Nil(t, err)
	req.Header.Set("Content-Type", wr.FormDataContentType())
	respWriter := httptest.NewRecorder()
	handler.EditStaffHandler(respWriter, req)
	resp := respWriter.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var responseStruct staffHttpResponse
	err = json.Unmarshal(body, &responseStruct)
	assert.NoError(t, err)
	assert.Equal(t, responseStruct.Data, returnStaff)

}

//
//func TestAdd(t *testing.T) {
//	//todo паника во время сохранения
//
//
//	type addstaffHandlerTestCase struct {
//		inputstaff  *staffModels.Staff
//		outputstaff staffModels.SafeStaff
//		httpErrs    []responses.HttpError
//	}
//
//	mockstaffUcase := new(mocks.Usecase)
//	handler := http2.StaffHandler{SUsecase: mockstaffUcase}
//	var inputstaff staffModels.Staff
//	err := faker.FakeData(&inputstaff)
//	assert.NoError(t, err)
//	inputstaff.EditedAt = time.Now().UTC()
//
//	outputstaff := app.GetSafeStaff(inputstaff)
//
//	inputstaff.StaffID = 0
//
//	testCases := []addstaffHandlerTestCase{
//		//Test OK
//		{
//			inputstaff:  &inputstaff,
//			outputstaff: outputstaff,
//			httpErrs:    nil,
//		},
//		//Test empty JsonData
//		{
//			inputstaff:  nil,
//			outputstaff: staffModels.SafeStaff{},
//			httpErrs: []responses.HttpError{{
//				Code:    400,
//				Message: globalModels.ErrEmptyJSON.Error(),
//			},
//			},
//		},
//	}
//
//	for i, testCase := range testCases {
//		message := fmt.Sprintf("test case number: %d", i)
//		mockstaffUcase.On("Add",
//			mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("models.Staff")).Return(testCase.outputstaff, nil)
//		mockstaffUcase.On("GetCafeId", mock.AnythingOfType("*context.valueCtx"), "").Return(2, nil)
//		mockstaffUcase.On("DeleteQrCodes", "").Return(nil)
//
//		var buf bytes.Buffer
//		var wr *multipart.Writer
//		if testCase.inputstaff != nil {
//			requestData, err := json.Marshal(&testCase.inputstaff)
//			assert.NoError(t, err, message)
//			buf, wr = createMultipartFormData(t, string(requestData))
//		} else {
//			buf, wr = createMultipartFormData(t, "")
//		}
//
//		req, err := http.NewRequest("POST", url, &buf)
//		session := sessions.Session{Values: map[interface{}]interface{}{"userID": testCase.inputstaff.StaffID}}
//		req = req.WithContext(context.WithValue(req.Context(), "session", &session))
//		assert.NoError(t, err, message)
//		req.Header.Set("Content-Type", wr.FormDataContentType())
//
//		respWriter := httptest.NewRecorder()
//
//		handler.AddStaffHandler(respWriter, req)
//
//		resp := respWriter.Result()
//		body, err := ioutil.ReadAll(resp.Body)
//		assert.NoError(t, err, message)
//
//		var responseStruct staffHttpResponse
//		err = json.Unmarshal(body, &responseStruct)
//		assert.NoError(t, err, message)
//
//		errs := responseStruct.Errors
//		staff := testCase.outputstaff
//
//		assert.Equal(t, testCase.httpErrs, errs, message)
//		assert.Equal(t, testCase.outputstaff, staff, message)
//	}
//
//}
//
