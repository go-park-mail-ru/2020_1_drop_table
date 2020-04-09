package test

import (
	http3 "2020_1_drop_table/internal/app/customer/delivery/http"
	"2020_1_drop_table/internal/app/customer/mocks"
	"2020_1_drop_table/internal/app/customer/models"
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

type custHttpResponse struct {
	Data   int
	Errors []responses.HttpError
}

func TestGet(t *testing.T) {
	const url = "/api/v1/asd/customer"
	mockcustomerUcase := new(mocks.Usecase)
	returncustomer := models.Customer{
		CustomerID: "",
		CafeID:     0,
		Points:     228,
	}
	mockcustomerUcase.On("GetPoints", mock.AnythingOfType("*context.valueCtx"), "asd").Return(returncustomer.Points, nil)
	handler := http3.CustomerHandler{CUsecase: mockcustomerUcase}
	buf, wr := createMultipartFormData(t, "")

	req, err := http.NewRequest("GET", url, &buf)
	req = mux.SetURLVars(req, map[string]string{
		"uuid": "asd",
	})
	assert.Nil(t, err)
	req.Header.Set("Content-Type", wr.FormDataContentType())
	respWriter := httptest.NewRecorder()

	handler.GetPoints(respWriter, req)
	resp := respWriter.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var responseStruct custHttpResponse
	err = json.Unmarshal(body, &responseStruct)
	assert.NoError(t, err)
	assert.Equal(t, responseStruct.Data, returncustomer.Points)
}

func TestSet(t *testing.T) {
	const url = "/api/v1/asd/customer/228"
	mockcustomerUcase := new(mocks.Usecase)
	returncustomer := models.Customer{
		CustomerID: "asd",
		CafeID:     0,
		Points:     228,
	}
	mockcustomerUcase.On("SetPoints", mock.AnythingOfType("*context.valueCtx"), "asd", 228).Return(nil)
	handler := http3.CustomerHandler{CUsecase: mockcustomerUcase}
	buf, wr := createMultipartFormData(t, "")

	req, err := http.NewRequest("PUT", url, &buf)
	req = mux.SetURLVars(req, map[string]string{
		"points": "228",
		"uuid":   "asd",
	})

	assert.Nil(t, err)
	req.Header.Set("Content-Type", wr.FormDataContentType())
	respWriter := httptest.NewRecorder()

	handler.SetPoints(respWriter, req)
	resp := respWriter.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var responseStruct custHttpResponse
	err = json.Unmarshal(body, &responseStruct)
	assert.NoError(t, err)
	assert.Equal(t, responseStruct.Data, returncustomer.Points)
}
