package cafes

import (
	"2020_1_drop_table/owners"
	"2020_1_drop_table/responses"
	"2020_1_drop_table/utils/testsUtils"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func CreateUserForTest(email, password string) (owners.Staff, error) {
	user := owners.Staff{
		Name:     "Василий Андреев",
		Email:    email,
		Password: password,
	}
	stf, err := owners.Storage.Append(user)

	return stf, err
}

type HttpTestCase struct {
	Context    map[string]string
	Cookie     http.Cookie
	Request    interface{}
	Response   responses.HttpResponse
	StatusCode int
}

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
	w.Close()
	return b, w
}

func TestCafeCreation(t *testing.T) {
	Storage.Clear()
	owners.Storage.Clear()

	email := "TestCafeCreation@example.com"
	password := "PassWord1"

	staff, err := CreateUserForTest(email, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	cafeOK := Cafe{
		Name:        "Пушкин",
		Address:     "Тверской б-р, 26А, Москва, 125009",
		Description: "Описание",
	}

	cafeNotOK := Cafe{
		Address:     "Тверской б-р, 26А, Москва, 125009",
		Description: "Описание",
	}
	testCases := []HttpTestCase{
		{
			Request: cafeOK,
			Response: responses.HttpResponse{
				Data:   cafeOK,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: cafeNotOK,
			Response: responses.HttpResponse{
				Data: cafeOK,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "Name is a required field",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
	}

	authCookieStaff1, err := testsUtils.GetAuthCookie(staff.StaffID)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}
	url := "/api/v1/cafe"

	for caseNum, item := range testCases {
		requestData, _ := json.Marshal(item.Request)
		var req *http.Request
		if requestData != nil {
			b, w := createMultipartFormData(t, string(requestData))
			req = httptest.NewRequest("POST", url, &b)
			req.Header.Set("Content-Type", w.FormDataContentType())
		} else {
			req = httptest.NewRequest("POST", url, nil)
		}

		respWriter := httptest.NewRecorder()

		req.AddCookie(&authCookieStaff1)

		CreateCafeHandler(respWriter, req)

		resp := respWriter.Result()
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var TrueResponse responses.HttpResponse

		err := json.Unmarshal(body, &TrueResponse)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		switch TrueResponse.Errors {
		case nil:
			//Data equals
			responseData := TrueResponse.Data.(map[string]interface{})
			expectedData := item.Response.Data.(Cafe)

			if responseData["name"] != expectedData.Name {
				t.Errorf("[%d] wrong Name field in response data: got %+v, expected %+v",
					caseNum, responseData["name"], expectedData.Name)
			}

			if responseData["address"] != expectedData.Address {
				t.Errorf("[%d] wrong Address field in response data: got %+v, expected %+v",
					caseNum, responseData["address"], expectedData.Address)
			}

			if responseData["description"] != expectedData.Description {
				t.Errorf("[%d] wrong Description field in response data: got %+v, expected %+v",
					caseNum, responseData["description"], expectedData.Description)
			}

		default:
			//Error equal
			if len(TrueResponse.Errors) != len(item.Response.Errors) {
				t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
					caseNum, len(TrueResponse.Errors), len(item.Response.Errors))
			}

			for errorNum, err := range TrueResponse.Errors {
				if err != item.Response.Errors[errorNum] {
					t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
						caseNum, err, item.Response.Errors[errorNum])
				}
			}
		}
	}
}

func createCafeForTest(cafeName string, staffID int) (Cafe, error) {
	cafe := Cafe{
		Name:        cafeName,
		Address:     "Тверской б-р, 26А, Москва, 125009",
		Description: "Описание",
		StaffID:     staffID,
	}
	return Storage.Append(cafe)
}

func TestGetCafeList(t *testing.T) {
	Storage.Clear()
	owners.Storage.Clear()

	email1 := "TestGetCafeList1@example.com"
	email2 := "TestGetCafeList2@example.com"
	password := "PassWord1"

	staff1, err := CreateUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieStaff1, err := testsUtils.GetAuthCookie(staff1.StaffID)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	staff2, err := CreateUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieStaff2, err := testsUtils.GetAuthCookie(staff2.StaffID)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	cafeName1 := "TestGetCafeList1"
	cafe1, err := createCafeForTest(cafeName1, staff1.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName2 := "TestGetCafeList2"
	cafe2, err := createCafeForTest(cafeName2, staff1.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName3 := "TestGetCafeList3"
	cafe3, err := createCafeForTest(cafeName3, staff2.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	testCases := []HttpTestCase{
		{
			Request: nil,
			Cookie:  authCookieStaff1,
			Response: responses.HttpResponse{
				Data: []Cafe{
					cafe1,
					cafe2,
				},
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: nil,
			Cookie:  authCookieStaff2,
			Response: responses.HttpResponse{
				Data: []Cafe{
					cafe3,
				},
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: nil,
			Cookie:  http.Cookie{},
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no permissions",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
	}

	url := "/api/v1/cafe"
	method := "GET"

	for caseNum, item := range testCases {
		req := httptest.NewRequest(method, url, nil)
		respWriter := httptest.NewRecorder()

		req.AddCookie(&item.Cookie)

		GetCafesListHandler(respWriter, req)

		resp := respWriter.Result()
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var TrueResponse responses.HttpResponse

		err = json.Unmarshal(body, &TrueResponse)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		switch TrueResponse.Errors {
		case nil:
			//Data equals
			trueResponse := TrueResponse.Data.([]interface{})
			expectedData := item.Response.Data.([]Cafe)

			if len(trueResponse) != len(expectedData) {
				t.Errorf("[%d] wrong cafe slice len: got %+v, expected %+v",
					caseNum, len(trueResponse), len(expectedData))
			}
			for i := range trueResponse {
				tmpResponse := trueResponse[i].(map[string]interface{})
				if tmpResponse["name"] != expectedData[i].Name {
					t.Errorf("[%d] wrong Name field in response data: got %+v, expected %+v",
						caseNum, tmpResponse["name"], expectedData[i].Name)
				}

				if tmpResponse["staffID"].(float64) != float64(expectedData[i].StaffID) {
					t.Errorf("[%d] wrong StaffID field in response data: got %+v, expected %+v",
						caseNum, tmpResponse["staffID"], expectedData[i].StaffID)
				}
			}

		default:
			//Error equal
			if len(TrueResponse.Errors) != len(item.Response.Errors) {
				t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
					caseNum, len(TrueResponse.Errors), len(item.Response.Errors))
			}

			for errorNum, err := range TrueResponse.Errors {
				if err != item.Response.Errors[errorNum] {
					t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
						caseNum, err, item.Response.Errors[errorNum])
				}
			}
		}

	}
}

func TestGetCafeHandler(t *testing.T) {
	Storage.Clear()
	owners.Storage.Clear()

	email1 := "TestGetCafeHandler1@example.com"
	email2 := "TestGetCafeHandler2@example.com"
	password := "PassWord1"

	staff1, err := CreateUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieStaff1, err := testsUtils.GetAuthCookie(staff1.StaffID)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	staff2, err := CreateUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieStaff2, err := testsUtils.GetAuthCookie(staff2.StaffID)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	cafeName1 := "TestGetCafeList1"
	cafe1, err := createCafeForTest(cafeName1, staff1.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName2 := "TestGetCafeList2"
	cafe2, err := createCafeForTest(cafeName2, staff2.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	testCases := []HttpTestCase{
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.CafeID)},
			Request: nil,
			Cookie:  authCookieStaff1,
			Response: responses.HttpResponse{
				Data:   cafe1,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe2.CafeID)},
			Request: nil,
			Cookie:  authCookieStaff2,
			Response: responses.HttpResponse{
				Data:   cafe2,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.CafeID)},
			Request: nil,
			Cookie:  authCookieStaff2,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no permissions",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe2.CafeID)},
			Request: nil,
			Cookie:  authCookieStaff1,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no permissions",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
	}

	url := "/api/v1/cafe"
	method := "GET"

	for caseNum, item := range testCases {
		req := httptest.NewRequest(method, url, nil)
		respWriter := httptest.NewRecorder()

		req.AddCookie(&item.Cookie)

		req = mux.SetURLVars(req, item.Context)

		GetCafeHandler(respWriter, req)

		resp := respWriter.Result()
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var TrueResponse responses.HttpResponse

		err = json.Unmarshal(body, &TrueResponse)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		switch TrueResponse.Errors {
		case nil:
			//Data equals
			trueResponse := TrueResponse.Data.(map[string]interface{})
			expectedData := item.Response.Data.(Cafe)

			if trueResponse["name"] != expectedData.Name {
				t.Errorf("[%d] wrong Name field in response data: got %+v, expected %+v",
					caseNum, trueResponse["name"], expectedData.Name)
			}

			if trueResponse["staffID"].(float64) != float64(expectedData.StaffID) {
				t.Errorf("[%d] wrong StaffID field in response data: got %+v, expected %+v",
					caseNum, trueResponse["staffID"], expectedData.StaffID)
			}

		default:
			//Error equal
			if len(TrueResponse.Errors) != len(item.Response.Errors) {
				t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
					caseNum, len(TrueResponse.Errors), len(item.Response.Errors))
			}

			for errorNum, err := range TrueResponse.Errors {
				if err != item.Response.Errors[errorNum] {
					t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
						caseNum, err, item.Response.Errors[errorNum])
				}
			}
		}
	}
}

func TestEditCafeHandler(t *testing.T) {
	Storage.Clear()
	owners.Storage.Clear()

	email1 := "TestEditCafeHandler1@example.com"
	email2 := "TestEditCafeHandler2@example.com"
	password := "PassWord1"

	staff1, err := CreateUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieStaff1, err := testsUtils.GetAuthCookie(staff1.StaffID)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	staff2, err := CreateUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieStaff2, err := testsUtils.GetAuthCookie(staff2.StaffID)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	_, err = createCafeForTest("TestGetCafeList1", staff1.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName1 := "TestGetCafeList2"
	cafe1, err := createCafeForTest(cafeName1, staff1.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}
	cafe1.Name = "TestGetCafeList1EDITED"

	cafeName2 := "TestGetCafeList2"
	cafe2, err := createCafeForTest(cafeName2, staff2.StaffID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}
	cafe2.Name = "TestGetCafeList2EDITED"

	testCases := []HttpTestCase{
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.CafeID)},
			Request: cafe1,
			Cookie:  authCookieStaff1,
			Response: responses.HttpResponse{
				Data:   cafe1,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe2.CafeID)},
			Request: cafe2,
			Cookie:  authCookieStaff2,
			Response: responses.HttpResponse{
				Data:   cafe2,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.CafeID)},
			Request: cafe1,
			Cookie:  authCookieStaff2,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no permissions",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.CafeID)},
			Request: cafe1,
			Cookie:  http.Cookie{},
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no permissions",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": "This is not CafeID"},
			Request: cafe1,
			Cookie:  authCookieStaff1,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "bad id: This is not CafeID",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": "1234567890"},
			Request: cafe1,
			Cookie:  authCookieStaff1,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no permissions",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.CafeID)},
			Request: Cafe{
				CafeID:  1,
				Name:    "Name",
				Address: "Address",
			},
			Cookie: authCookieStaff1,
			Response: responses.HttpResponse{
				Data: cafe1,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "Description is a required field",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
	}
	url := "/api/v1/cafe"
	method := "PUT"

	for caseNum, item := range testCases {
		requestData, _ := json.Marshal(item.Request)
		var req *http.Request
		if requestData != nil {
			b, w := createMultipartFormData(t, string(requestData))
			req = httptest.NewRequest(method, url, &b)
			req.Header.Set("Content-Type", w.FormDataContentType())
		} else {
			req = httptest.NewRequest(method, url, nil)
		}

		respWriter := httptest.NewRecorder()

		req.AddCookie(&item.Cookie)

		req = mux.SetURLVars(req, item.Context)

		EditCafeHandler(respWriter, req)

		resp := respWriter.Result()
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var TrueResponse responses.HttpResponse

		err = json.Unmarshal(body, &TrueResponse)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		switch TrueResponse.Errors {
		case nil:
			//Data equals
			responseData := TrueResponse.Data.(map[string]interface{})
			expectedData := item.Response.Data.(Cafe)

			if responseData["name"] != expectedData.Name {
				t.Errorf("[%d] wrong Name field in response data: got %+v, expected %+v",
					caseNum, responseData["name"], expectedData.Name)
			}

			if responseData["address"] != expectedData.Address {
				t.Errorf("[%d] wrong Address field in response data: got %+v, expected %+v",
					caseNum, responseData["address"], expectedData.Address)
			}

			if responseData["description"] != expectedData.Description {
				t.Errorf("[%d] wrong Description field in response data: got %+v, expected %+v",
					caseNum, responseData["description"], expectedData.Description)
			}

		default:
			//Error equal
			if len(TrueResponse.Errors) != len(item.Response.Errors) {
				t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
					caseNum, len(TrueResponse.Errors), len(item.Response.Errors))
			}

			for errorNum, err := range TrueResponse.Errors {
				if err != item.Response.Errors[errorNum] {
					t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
						caseNum, err, item.Response.Errors[errorNum])
				}
			}
		}
	}
}
