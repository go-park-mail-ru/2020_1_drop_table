package middleware

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
