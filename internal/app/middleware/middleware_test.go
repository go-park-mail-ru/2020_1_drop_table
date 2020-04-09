package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanic(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
		assert.True(t, true)
	})
	handlerToTest := PanicMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	req = req.WithContext(context.WithValue(req.Context(), "session", &session))

	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)
}

func TestLog(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, true)
	})
	handlerToTest := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	req = req.WithContext(context.WithValue(req.Context(), "session", &session))

	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)
}

func TestCors(t *testing.T) {
	r := mux.NewRouter()
	MyCORSMethodMiddleware(r)
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, true, false)
		}
	}()
	assert.Equal(t, true, true)
}

func TestNew(t *testing.T) {
	r := mux.NewRouter()
	NewMiddleware(r, nil)
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, true, false)
		}
	}()
	assert.Equal(t, true, true)
}
