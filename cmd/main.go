package main

import (
	"2020_1_drop_table/configs"
	_cafeHttpDeliver "2020_1_drop_table/internal/app/cafe/delivery/http"
	_cafeRepo "2020_1_drop_table/internal/app/cafe/repository"
	_cafeUsecase "2020_1_drop_table/internal/app/cafe/usecase"
	"2020_1_drop_table/internal/app/middleware"
	_staffHttpDeliver "2020_1_drop_table/internal/app/staff/delivery/http"
	_staffRepo "2020_1_drop_table/internal/app/staff/repository"
	_staffUsecase "2020_1_drop_table/internal/app/staff/usecase"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	redisStore "gopkg.in/boj/redistore.v1"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()

	//Middleware
	var CookieStore, err = redisStore.NewRediStore(
		configs.RedisPreferences.Size,
		configs.RedisPreferences.Network,
		configs.RedisPreferences.Address,
		configs.RedisPreferences.Password,
		configs.RedisPreferences.SecretKey)

	middleware.NewMiddleware(r, CookieStore)

	timeoutContext := time.Second * 2
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s",
		configs.PostgresPreferences.User,
		configs.PostgresPreferences.Password,
		configs.PostgresPreferences.Port)

	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Error().Msgf(err.Error())
	}

	staffRepo := _staffRepo.NewPostgresStaffRepository(conn)

	if err != nil {
		log.Error().Msgf(err.Error())
	}
	cafeRepo := _cafeRepo.NewPostgresCafeRepository(conn)
	staffUsecase := _staffUsecase.NewStaffUsecase(&staffRepo, &cafeRepo, timeoutContext)
	_staffHttpDeliver.NewStaffHandler(r, staffUsecase)

	if err != nil {
		log.Error().Msgf(err.Error())
	}
	cafeUsecase := _cafeUsecase.NewCafeUsecase(&cafeRepo, staffUsecase, timeoutContext)
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
		http.StripPrefix("/media/", http.FileServer(http.Dir(configs.MediaFolder))))

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
