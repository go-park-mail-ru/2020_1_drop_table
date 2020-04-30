package middleware

import (
	"2020_1_drop_table/internal/pkg/metrics"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

func NewPanicMiddleware(metrics *metrics.PromMetrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqTime := time.Now()
			defer func() {
				if err := recover(); err != nil {
					respTime := time.Since(reqTime)
					metrics.Hits.WithLabelValues(
						strconv.Itoa(http.StatusInternalServerError), r.URL.Path, r.Method).Inc()

					metrics.Timings.WithLabelValues(
						strconv.Itoa(http.StatusInternalServerError), r.URL.String(),
						r.Method).Observe(respTime.Seconds())

					log.Error().Msgf(fmt.Sprintf("panic catched: %s", err))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
