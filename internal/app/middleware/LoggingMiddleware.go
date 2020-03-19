package middleware

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("URL: %s, METHOD: %s", r.RequestURI, r.Method)
		log.Info().Msgf(msg)
		next.ServeHTTP(w, r)
	})
}
