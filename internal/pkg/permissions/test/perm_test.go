package test

import (
	"2020_1_drop_table/internal/pkg/permissions"
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetCsrf(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrf := w.Header().Get("csrf")

		assert.NotEqual(t, "", csrf)
	})
	handlerToTest := permissions.SetCSRF(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)

	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)

}

func TestCheckCsrf(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrf := w.Header().Get("csrf")
		fmt.Println(csrf)
		assert.Equal(t, "", csrf)
	})
	handlerToTest := permissions.CheckCSRF(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)

	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)
}

func TestCheckAuth(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, true)
	})
	handlerToTest := permissions.CheckAuthenticated(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	req = req.WithContext(context.WithValue(req.Context(), "session", &session))

	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)
}
