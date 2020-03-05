package main

import (
	"2020_1_drop_table/cafes"
	"2020_1_drop_table/middlewares"
	"2020_1_drop_table/owners"
	"2020_1_drop_table/projectConfig"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()

	//Middleware
	r.Use(middlewares.MyCORSMethodMiddleware(r))
	r.Use(middlewares.LoggingMiddleware)

	//owner handlers
	r.HandleFunc("/api/v1/owner", owners.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/owner/login", owners.LoginHandler).Methods("POST")
	r.HandleFunc("/api/v1/owner/{id:[0-9]+}", owners.GetOwnerHandler).Methods("GET")
	r.HandleFunc("/api/v1/getCurrentOwner/", owners.GetCurrentOwnerHandler).Methods("GET")
	r.HandleFunc("/api/v1/owner/{id:[0-9]+}", owners.EditOwnerHandler).Methods("PUT")

	//cafe handlers
	r.HandleFunc("/api/v1/cafe", cafes.CreateCafeHandler).Methods("POST")
	r.HandleFunc("/api/v1/cafe", cafes.GetCafesListHandler).Methods("GET")
	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}", cafes.GetCafeHandler).Methods("GET")
	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}", cafes.EditCafeHandler).Methods("PUT")

	//OPTIONS
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length,"+
			" Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers,"+
			" Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	//static server
	r.PathPrefix("/media/").Handler(
		http.StripPrefix("/media/", http.FileServer(http.Dir(projectConfig.MediaFolder))))

	http.Handle("/", r)
	log.Info().Msgf("starting server at :8080")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Error().Msgf(srv.ListenAndServe().Error())
}
