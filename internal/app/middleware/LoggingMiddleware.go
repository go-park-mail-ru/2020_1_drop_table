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

func NewLoggingMiddleware(metrics *metrics.PromMetrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI != "/metrics" {
				msg := fmt.Sprintf("URL: %s, METHOD: %s", r.RequestURI, r.Method)
				log.Info().Msgf(msg)
			}

			reqTime := time.Now()
			next.ServeHTTP(w, r)
			respTime := time.Since(reqTime)
			if r.URL.Path != "/metrics" {
				//If status is 500 => response will be in panic middleware
				metrics.Hits.WithLabelValues(strconv.Itoa(http.StatusOK), r.URL.String(), r.Method).Inc()
				metrics.Timings.WithLabelValues(
					strconv.Itoa(http.StatusOK), r.URL.String(), r.Method).Observe(respTime.Seconds())
			}
		})
	}
}
