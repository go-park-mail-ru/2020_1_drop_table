package main

import (
	"2020_1_drop_table/cafes"
	"2020_1_drop_table/owners"
	"2020_1_drop_table/responses"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCafeCreation(t *testing.T) {
	email := "TestCafeCreation@example.com"
	password := "PassWord1"

	err, _ := createUserForTest(email, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	cafeOK := cafes.Cafe{
		Name:        "Пушкин",
		Address:     "Тверской б-р, 26А, Москва, 125009",
		Description: "Описание",
	}

	cafeNotOK := cafes.Cafe{
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

	authCookieOwner1, err := owners.GetAuthCookie(email, password)
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

		req.AddCookie(&authCookieOwner1)

		cafes.CreateCafeHandler(respWriter, req)

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
			expectedData := item.Response.Data.(cafes.Cafe)

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

func createCafeForTest(cafeName string, ownerID int) (error, cafes.Cafe) {
	cafe := cafes.Cafe{
		Name:        cafeName,
		Address:     "Тверской б-р, 26А, Москва, 125009",
		Description: "Описание",
		OwnerID:     ownerID,
	}
	return cafes.Storage.Append(cafe)
}

func TestGetCafeList(t *testing.T) {
	email1 := "TestGetCafeList1@example.com"
	email2 := "TestGetCafeList2@example.com"
	password := "PassWord1"

	err, owner1 := createUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieOwner1, err := owners.GetAuthCookie(email1, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	err, owner2 := createUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieOwner2, err := owners.GetAuthCookie(email2, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	cafeName1 := "TestGetCafeList1"
	err, cafe1 := createCafeForTest(cafeName1, owner1.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName2 := "TestGetCafeList2"
	err, cafe2 := createCafeForTest(cafeName2, owner1.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName3 := "TestGetCafeList3"
	err, cafe3 := createCafeForTest(cafeName3, owner2.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	testCases := []HttpTestCase{
		{
			Request: nil,
			Cookie:  authCookieOwner1,
			Response: responses.HttpResponse{
				Data: []cafes.Cafe{
					cafe1,
					cafe2,
				},
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Request: nil,
			Cookie:  authCookieOwner2,
			Response: responses.HttpResponse{
				Data: []cafes.Cafe{
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

		cafes.GetCafesListHandler(respWriter, req)

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
			expectedData := item.Response.Data.([]cafes.Cafe)

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

				if tmpResponse["ownerID"].(float64) != float64(expectedData[i].OwnerID) {
					t.Errorf("[%d] wrong OwnerID field in response data: got %+v, expected %+v",
						caseNum, tmpResponse["ownerID"], expectedData[i].OwnerID)
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
	email1 := "TestGetCafeHandler1@example.com"
	email2 := "TestGetCafeHandler2@example.com"
	password := "PassWord1"

	err, owner1 := createUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieOwner1, err := owners.GetAuthCookie(email1, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	err, owner2 := createUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieOwner2, err := owners.GetAuthCookie(email2, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	cafeName1 := "TestGetCafeList1"
	err, cafe1 := createCafeForTest(cafeName1, owner1.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName2 := "TestGetCafeList2"
	err, cafe2 := createCafeForTest(cafeName2, owner2.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	testCases := []HttpTestCase{
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.ID)},
			Request: nil,
			Cookie:  authCookieOwner1,
			Response: responses.HttpResponse{
				Data:   cafe1,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe2.ID)},
			Request: nil,
			Cookie:  authCookieOwner2,
			Response: responses.HttpResponse{
				Data:   cafe2,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.ID)},
			Request: nil,
			Cookie:  authCookieOwner2,
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
			Context: map[string]string{"id": strconv.Itoa(cafe2.ID)},
			Request: nil,
			Cookie:  authCookieOwner1,
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

		cafes.GetCafeHandler(respWriter, req)

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
			expectedData := item.Response.Data.(cafes.Cafe)

			if trueResponse["name"] != expectedData.Name {
				t.Errorf("[%d] wrong Name field in response data: got %+v, expected %+v",
					caseNum, trueResponse["name"], expectedData.Name)
			}

			if trueResponse["ownerID"].(float64) != float64(expectedData.OwnerID) {
				t.Errorf("[%d] wrong OwnerID field in response data: got %+v, expected %+v",
					caseNum, trueResponse["ownerID"], expectedData.OwnerID)
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
	email1 := "TestEditCafeHandler1@example.com"
	email2 := "TestEditCafeHandler2@example.com"
	password := "PassWord1"

	err, owner1 := createUserForTest(email1, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieOwner1, err := owners.GetAuthCookie(email1, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	err, owner2 := createUserForTest(email2, password)
	if err != nil {
		t.Errorf("can't create new user, error: %+v", err)
	}

	authCookieOwner2, err := owners.GetAuthCookie(email2, password)
	if err != nil {
		t.Errorf("auth error: %s", err)
	}

	err, _ = createCafeForTest("TestGetCafeList1", owner1.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}

	cafeName1 := "TestGetCafeList2"
	err, cafe1 := createCafeForTest(cafeName1, owner1.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}
	cafe1.Name = "TestGetCafeList1EDITED"

	cafeName2 := "TestGetCafeList2"
	err, cafe2 := createCafeForTest(cafeName2, owner2.ID)
	if err != nil {
		t.Errorf("got error while creationg cafe: %+v", err)
	}
	cafe2.Name = "TestGetCafeList2EDITED"

	testCases := []HttpTestCase{
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.ID)},
			Request: cafe1,
			Cookie:  authCookieOwner1,
			Response: responses.HttpResponse{
				Data:   cafe1,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe2.ID)},
			Request: cafe2,
			Cookie:  authCookieOwner2,
			Response: responses.HttpResponse{
				Data:   cafe2,
				Errors: nil,
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": strconv.Itoa(cafe1.ID)},
			Request: cafe1,
			Cookie:  authCookieOwner2,
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
			Context: map[string]string{"id": strconv.Itoa(cafe1.ID)},
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
			Context: map[string]string{"id": "This is not ID"},
			Request: cafe1,
			Cookie:  authCookieOwner1,
			Response: responses.HttpResponse{
				Data: nil,
				Errors: []responses.HttpError{
					{
						Code:    400,
						Message: "bad id: This is not ID",
					},
				},
			},
			StatusCode: http.StatusOK,
		},
		{
			Context: map[string]string{"id": "1234567890"},
			Request: cafe1,
			Cookie:  authCookieOwner1,
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
			Context: map[string]string{"id": strconv.Itoa(cafe1.ID)},
			Request: cafes.Cafe{
				ID:      1,
				Name:    "Name",
				Address: "Address",
			},
			Cookie: authCookieOwner1,
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

		cafes.EditCafeHandler(respWriter, req)

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
			expectedData := item.Response.Data.(cafes.Cafe)

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
