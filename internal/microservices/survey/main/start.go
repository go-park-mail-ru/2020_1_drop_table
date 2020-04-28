package main

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app/middleware"
	http2 "2020_1_drop_table/internal/microservices/survey/delivery/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	timeoutContext := configs.Timeouts.ContextTimeout

	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s",
		configs.PostgresPreferences.User,
		configs.PostgresPreferences.Password,
		configs.PostgresPreferences.Port)

	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Error().Msgf(err.Error())
		return
	}

	survRepo := surveyRepo.NewPostgresSurveyRepository(conn)
	surveyUcase := surveyUsecase.NewSurveyUsecase(cafeRepo, survRepo, staffUsecase, timeoutContext)
	http2.NewSurveyHandler(r, surveyUcase)

	//OPTIONS
	middleware.AddOptionsRequest(r)

	//static server
	r.PathPrefix(fmt.Sprintf("/%s/", configs.MediaFolder)).Handler(
		http.StripPrefix(fmt.Sprintf("/%s/", configs.MediaFolder),
			http.FileServer(http.Dir(configs.MediaFolder))))

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: configs.Timeouts.WriteTimeout,
		ReadTimeout:  configs.Timeouts.ReadTimeout,
	}
	log.Error().Msgf(srv.ListenAndServe().Error())
