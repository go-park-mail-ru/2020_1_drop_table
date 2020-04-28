package main

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app/middleware"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/test_client"
	http2 "2020_1_drop_table/internal/microservices/survey/delivery/http"
	surveyRepo "2020_1_drop_table/internal/microservices/survey/repository"
	surveyUsecase "2020_1_drop_table/internal/microservices/survey/usecase"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	redisStore "gopkg.in/boj/redistore.v1"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	var CookieStore, err = redisStore.NewRediStore(
		configs.RedisPreferences.Size,
		configs.RedisPreferences.Network,
		configs.RedisPreferences.Address,
		configs.RedisPreferences.Password,
		configs.RedisPreferences.SecretKey)

	middleware.NewMiddleware(r, CookieStore)

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

	grpcConn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	grpcStaffClient := staffClient.NewStaffClient(grpcConn)
	surveyUcase := surveyUsecase.NewSurveyUsecase(survRepo, grpcStaffClient, timeoutContext)
	http2.NewSurveyHandler(r, surveyUcase)

	//OPTIONS
	middleware.AddOptionsRequest(r)

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8085",
		WriteTimeout: configs.Timeouts.WriteTimeout,
		ReadTimeout:  configs.Timeouts.ReadTimeout,
	}
	log.Error().Msgf(srv.ListenAndServe().Error())
}
