package middleware

import (
	"2020_1_drop_table/configs"
	"github.com/gorilla/mux"
	"net/http"
)

func MyCORSMethodMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods",
				"POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length,"+
				" Accept-Encoding, X-CSRF-Token, csrf-token, Authorization")
			w.Header().Set("Access-Control-Allow-Origin", configs.FrontEndUrl)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Accept, Cookie")
			if req.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
