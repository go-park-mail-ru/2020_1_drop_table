package middleware

import (
	"2020_1_drop_table/internal/pkg/metrics"
	"github.com/gorilla/mux"
	"gopkg.in/boj/redistore.v1"
)

func NewMiddleware(r *mux.Router, sp *redistore.RediStore, metrics *metrics.PromMetrics) {
	sessionMiddleware := sessionMiddleware{
		sessionRepo: sp,
	}
	r.Use(MyCORSMethodMiddleware(r))
	r.Use(sessionMiddleware.SessionMiddleware)
	r.Use(NewLoggingMiddleware(metrics))
	r.Use(NewPanicMiddleware(metrics))
}
