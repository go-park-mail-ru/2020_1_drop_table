package main

import (
	"2020_1_drop_table/owners"
	"2020_1_drop_table/responses"
	"bytes"
	"encoding/json"
	"fmt"
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

func TestRegisterUser(t *testing.T) {
	ownerObjOK := owners.Owner{
		Name:     "Василий Андреев",
		Email:    "example@example.com",
		Password: "PassWord1",
	}

	ownerObjNotOK := owners.Owner{
		Name:  "Василий Андреев",
		Email: "example@example.com",
	}

	testCases := []HttpTestCase{
		{
			Request: ownerObjOK,
			Response: responses.HttpResponse{
				Data:   ownerObjOK,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: ownerObjNotOK,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    http.StatusBadRequest,
						Message: "Password is a required field",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: nil,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    http.StatusBadRequest,
						Message: "empty jsonData field",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: ownerObjOK,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    http.StatusBadRequest,
						Message: "User with this email already existed",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
	}

	url := "/api/v1/owner"
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

		owners.RegisterHandler(respWriter, req)

		resp := respWriter.Result()
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var responseObject responses.HttpResponse

		err := json.Unmarshal(body, &responseObject)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		if responseObject.Data != nil {
			//Data equals
			responseData := responseObject.Data.(map[string]interface{})
			expectedData := item.Response.Data.(owners.Owner)

			if responseData["name"] != expectedData.Name {
				t.Errorf("[%d] wrong Name field in response data: got %+v, expected %+v",
					caseNum, responseData["name"], expectedData.Name)
			}

			if responseData["email"] != expectedData.Email {
				t.Errorf("[%d] wrong Email field in response data: got %+v, expected %+v",
					caseNum, responseData["email"], expectedData.Email)
			}

		} else if item.Response.Errors != nil {
			//Error equal
			if len(responseObject.Errors) != len(item.Response.Errors) {
				t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
					caseNum, len(responseObject.Errors), len(item.Response.Errors))
			}

			for errorNum, err := range responseObject.Errors {
				if err != item.Response.Errors[errorNum] {
					t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
						caseNum, err, item.Response.Errors[errorNum])
				}
			}

		} else {
			t.Errorf("[%d] wrong response data: got nil, expected %+v",
				caseNum, item.Response.Data)
		}
	}
}

func CreateUserForTest(email, password string) (error, owners.Owner) {
	user := owners.Owner{
		Name:     "Василий Андреев",
		Email:    email,
		Password: password,
	}
	own, err := owners.Storage.Append(user)

	return err, own
}

func TestLoginUser(t *testing.T) {
	//Preparing for test
	email := "testLoginUser@example.com"
	password := "PassWord1"
	err, _ := CreateUserForTest(email, password)

	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}
	//Test
	testCases := []HttpTestCase{
		{
			Request: fmt.Sprintf(`{"email":  "%s",  "password": "%s"}`,
				email, password),
			Response: responses.HttpResponse{
				Data:   "",
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: fmt.Sprintf(`{"email":  "%s",  "password": "%sWrongPassword"}`,
				email, password),
			Response: responses.HttpResponse{
				Data: "",
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no user with given login and password",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: fmt.Sprintf(`{"email":  "%ssWrongEmail",  "password": "%s"}`, email, password),
			Response: responses.HttpResponse{
				Data: "",
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no user with given login and password",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: fmt.Sprintf(`{"email":  "%ssWrongEmail",  "password": "%sWrongPassword"}`, email, password),
			Response: responses.HttpResponse{
				Data: "",
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no user with given login and password",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: fmt.Sprintf(`{"email":  "%ssWrongEmail",  "password": "%sWrongPassword"}`, email, password),
			Response: responses.HttpResponse{
				Data: "",
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "no user with given login and password",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: fmt.Sprintf(`{"email":  "%ssWrongEmail"}`, email),
			Response: responses.HttpResponse{
				Data: "",
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "Password is a required field",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
	}
	url := "/api/v1/owner/login"
	for caseNum, item := range testCases {
		reader := strings.NewReader(item.Request.(string))
		req := httptest.NewRequest("POST", url, reader)

		respWriter := httptest.NewRecorder()
		owners.LoginHandler(respWriter, req)
		resp := respWriter.Result()

		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var trueResponse responses.HttpResponse
		err := json.Unmarshal(body, &trueResponse)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		if len(trueResponse.Errors) != len(item.Response.Errors) {
			t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
				caseNum, len(trueResponse.Errors), len(item.Response.Errors))
		}

		for errorNum, err := range trueResponse.Errors {
			if err != item.Response.Errors[errorNum] {
				t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
					caseNum, err, item.Response.Errors[errorNum])
			}
		}

		if len(trueResponse.Errors) == 0 {
			cookies := resp.Cookies()
			for _, cookie := range cookies {

				//Add new statement if new COOKIE will be added
				switch cookie.Name {
				case "authCookie":
					ownerFromCookie, err := owners.StorageSession.GetOwnerByCookie(cookie.Value)

					if err != nil {
						t.Errorf("[%d] error while getting error by Cookie: %+v:", caseNum, err)
					}

					if ownerFromCookie.Email != email {
						t.Errorf("[%d] wrong owner's email from COOKIE: got %+v, expected %+v",
							caseNum, ownerFromCookie.Email, email)
					}
				default:
					t.Errorf("[%d] unexpected Cookie with name: %+v:", caseNum, cookie.Name)
				}
			}
		}
	}
}

func TestGetOwner(t *testing.T) {
	//Preparing for test
	email1 := "testGetOwner1@example.com"
	email2 := "testGetOwner2@example.com"
	password := "PassWord1"

	err, owner1 := CreateUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	err, owner2 := CreateUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}
	//Test
	testCases := []HttpTestCase{
		{
			Context: map[string]string{"id": strconv.Itoa(owner2.OwnerID)},
			Request: nil,
			Response: responses.HttpResponse{
				Data: owners.Owner{
					OwnerID: owner2.OwnerID,
					Email:   owner2.Email,
				},
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(owner1.OwnerID)},
			Request: nil,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    http.StatusBadRequest,
						Message: "no permissions",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": "I'm not int"},
			Request: nil,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    http.StatusBadRequest,
						Message: "bad id: I'm not int",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
	}

	authCookieOwner1, err := owners.GetAuthCookie(email2, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}
	url := "/api/v1/owner"

	for caseNum, item := range testCases {
		req := httptest.NewRequest("GET", url, nil)
		respWriter := httptest.NewRecorder()

		req.AddCookie(&authCookieOwner1)

		req = mux.SetURLVars(req, item.Context)

		owners.GetOwnerHandler(respWriter, req)

		resp := respWriter.Result()
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var responseObject responses.HttpResponse

		err := json.Unmarshal(body, &responseObject)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		if len(item.Response.Errors) != len(responseObject.Errors) {
			t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
				caseNum, len(responseObject.Errors), len(item.Response.Errors))
		}

		switch responseObject.Errors {
		case nil:
			responseData := responseObject.Data.(map[string]interface{})
			expectedData := item.Response.Data.(owners.Owner)

			if responseData["id"].(float64) != float64(expectedData.OwnerID) {
				t.Errorf("[%d] wrong Name field in response data: got %+v, expected %+v",
					caseNum, responseData["id"], expectedData.OwnerID)
			}

			if responseData["email"] != expectedData.Email {
				t.Errorf("[%d] wrong Email field in response data: got %+v, expected %+v",
					caseNum, responseData["email"], expectedData.Email)
			}
		default:
			for errorNum, err := range responseObject.Errors {
				if err != item.Response.Errors[errorNum] {
					t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
						caseNum, err, item.Response.Errors[errorNum])
				}
			}
		}
	}
}

func TestGetCurrentOwner(t *testing.T) {
	//Preparing for test
	email1 := "GetGCurrentOwner1@example.com"
	email2 := "GetCurrentOwner2@example.com"
	password := "PassWord1"
	owners.Storage.Clear()

	err, owner1 := CreateUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	err, _ = CreateUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}
	//Test
	testCases := []HttpTestCase{
		{
			Request: nil,
			Response: responses.HttpResponse{
				Data:   owner1,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
	}

	authCookieOwner1, err := owners.GetAuthCookie(email1, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}
	url := "/api/v1/getCurrentOwner/"

	for caseNum, item := range testCases {
		req := httptest.NewRequest("GET", url, nil)
		respWriter := httptest.NewRecorder()

		req.AddCookie(&authCookieOwner1)

		owners.GetCurrentOwnerHandler(respWriter, req)

		resp := respWriter.Result()
		if resp.StatusCode != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %+v, expected %+v",
				caseNum, resp.StatusCode, item.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var responseObject responses.HttpResponse

		err := json.Unmarshal(body, &responseObject)
		if err != nil {
			t.Errorf("[%d] unmarshaling error: %s", caseNum, err)
		}

		if len(item.Response.Errors) != len(responseObject.Errors) {
			t.Errorf("[%d] wrong errors count in response: got %d, expected %d",
				caseNum, len(responseObject.Errors), len(item.Response.Errors))
		}

		switch responseObject.Errors {
		case nil:
			responseData := responseObject.Data.(map[string]interface{})
			expectedData := item.Response.Data.(owners.Owner)
			if responseData["id"].(float64) != float64(expectedData.OwnerID) {
				t.Errorf("[%d] wrong id field in response data: got %+v, expected %+v",
					caseNum, responseData["id"], expectedData.OwnerID)
			}

			if responseData["email"] != expectedData.Email {
				t.Errorf("[%d] wrong Email field in response data: got %+v, expected %+v",
					caseNum, responseData["email"], expectedData.Email)
			}
		default:
			for errorNum, err := range responseObject.Errors {
				if err != item.Response.Errors[errorNum] {
					t.Errorf("[%d] wrong error in response: got %+v, expected %+v",
						caseNum, err, item.Response.Errors[errorNum])
				}
			}
		}

	}
}

func TestEditOwnerHandler(t *testing.T) {
	//Preparing for test
	email1 := "testEditOwner1@example.com"
	email2 := "testEditOwner2@example.com"
	password := "PassWord1"
	owners.Storage.Clear()
	err, owner1 := CreateUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	err, owner2 := CreateUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}
	owner2.Email = "EDITED@EMAIL.com"

	authCookieOwner2, err := owners.GetAuthCookie(email2, password)
	//Test
	testCases := []HttpTestCase{
		{
			Cookie:  authCookieOwner2,
			Context: map[string]string{"id": strconv.Itoa(owner2.OwnerID)},
			Request: owner2,
			Response: responses.HttpResponse{
				Data: owners.Owner{
					OwnerID: owner2.OwnerID,
					Email:   owner2.Email,
				},
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Cookie:  authCookieOwner2,
			Context: map[string]string{"id": strconv.Itoa(owner1.OwnerID)},
			Request: nil,
			Response: responses.HttpResponse{
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
			Cookie:  http.Cookie{},
			Context: map[string]string{"id": "123456757"},
			Request: nil,
			Response: responses.HttpResponse{
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

	if err != nil {
		t.Errorf("auth error: %s", err)
	}
	url := "/api/v1/owner"

	for caseNum, item := range testCases {
		requestData, _ := json.Marshal(item.Request)
		var req *http.Request
		if requestData != nil {
			b, w := createMultipartFormData(t, string(requestData))
			req = httptest.NewRequest("PUT", url, &b)
			req.Header.Set("Content-Type", w.FormDataContentType())
		} else {
			req = httptest.NewRequest("PUT", url, nil)
		}

		respWriter := httptest.NewRecorder()

		req.AddCookie(&item.Cookie)

		req = mux.SetURLVars(req, item.Context)

		owners.EditOwnerHandler(respWriter, req)

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
			expectedData := item.Response.Data.(owners.Owner)

			if responseData["id"].(float64) != float64(expectedData.OwnerID) {
				t.Errorf("[%d] wrong CafeID field in response data: got %+v, expected %+v",
					caseNum, responseData["id"], expectedData.OwnerID)
			}

			if responseData["email"] != expectedData.Email {
				t.Errorf("[%d] wrong Email field in response data: got %+v, expected %+v",
					caseNum, responseData["email"], expectedData.Email)
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
