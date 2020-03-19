package main

import (
	_cafeHttpDeliver "2020_1_drop_table/app/cafe/delivery/http"
	_cafeRepo "2020_1_drop_table/app/cafe/repository"
	_cafeUsecase "2020_1_drop_table/app/cafe/usecase"
	"2020_1_drop_table/app/middleware"
	_staffHttpDeliver "2020_1_drop_table/app/staff/delivery/http"
	_staffRepo "2020_1_drop_table/app/staff/repository"
	_staffUsecase "2020_1_drop_table/app/staff/usecase"
	"2020_1_drop_table/projectConfig"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	redisStore "gopkg.in/boj/redistore.v1"
	"net/http"
	"os"
	"time"
)

func main() {
	r := mux.NewRouter()

	//Middleware
	var CookieStore, err = redisStore.NewRediStore(10, "tcp", ":6379",
		"", []byte(os.Getenv("SESSION_KEY")))

	middleware.NewMiddleware(r, CookieStore)

	timeoutContext := time.Second * 2
	//ToDo make file with project preferences
	staffRepo, err := _staffRepo.NewPostgresStaffRepository("postgres", "", "5431")

	if err != nil {
		log.Error().Msgf(err.Error())
	}
	staffUsecase := _staffUsecase.NewStaffUsecase(&staffRepo, timeoutContext)
	_staffHttpDeliver.NewStaffHandler(r, staffUsecase)

	cafeRepo, err := _cafeRepo.NewPostgresStaffRepository("postgres", "", "5431")
	if err != nil {
		log.Error().Msgf(err.Error())
	}
	cafeUsecase := _cafeUsecase.NewCafeUsecase(&cafeRepo, timeoutContext)
	_cafeHttpDeliver.NewCafeHandler(r, cafeUsecase)

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
