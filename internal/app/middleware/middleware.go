package middleware

import (
	"github.com/gorilla/mux"
	"gopkg.in/boj/redistore.v1"
)

func NewMiddleware(r *mux.Router, sp *redistore.RediStore) {
	sessionMiddleware := sessionMiddleware{
		sessionRepo: sp,
	}
	r.Use(MyCORSMethodMiddleware(r))
	r.Use(sessionMiddleware.SessionMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(PanicMiddleware)
}
