package main

import (
	"2020_1_drop_table/configs"
	cafeClient "2020_1_drop_table/internal/app/cafe/delivery/grpc/client"
	"2020_1_drop_table/internal/app/middleware"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	http2 "2020_1_drop_table/internal/microservices/survey/delivery/http"
	surveyRepo "2020_1_drop_table/internal/microservices/survey/repository"
	surveyUsecase "2020_1_drop_table/internal/microservices/survey/usecase"
	"2020_1_drop_table/internal/pkg/metrics"
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

	//PromMetrics server
	metricsProm := metrics.RegisterMetrics(r)

	//Middleware
	var CookieStore, err = redisStore.NewRediStore(
		configs.RedisPreferences.Size,
		configs.RedisPreferences.Network,
		configs.RedisPreferences.Address,
		configs.RedisPreferences.Password,
		configs.RedisPreferences.SecretKey)

	middleware.NewMiddleware(r, CookieStore, metricsProm)

	timeoutContext := configs.Timeouts.ContextTimeout

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%s host=%s",
		configs.PostgresPreferences.User,
		configs.PostgresPreferences.Password,
		configs.PostgresPreferences.DBName,
		configs.PostgresPreferences.Port,
		configs.PostgresPreferences.Host)

	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Error().Msgf(err.Error())
		return
	}

	survRepo := surveyRepo.NewPostgresSurveyRepository(conn)

	grpcConn, _ := grpc.Dial(configs.GRPCStaffUrl, grpc.WithInsecure())
	grpcStaffClient := staffClient.NewStaffClient(grpcConn)

	grpcCafeConn, _ := grpc.Dial(configs.GRPCCafeUrl, grpc.WithInsecure())
	grpcCafeClient := cafeClient.NewCafeClient(grpcCafeConn)

	surveyUcase := surveyUsecase.NewSurveyUsecase(survRepo, grpcStaffClient, grpcCafeClient, timeoutContext)
	http2.NewSurveyHandler(r, surveyUcase)

	//OPTIONS
	middleware.AddOptionsRequest(r)

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         configs.HTTPSurveyUrl,
		WriteTimeout: configs.Timeouts.WriteTimeout,
		ReadTimeout:  configs.Timeouts.ReadTimeout,
	}
	fmt.Println("survey server started at ", configs.HTTPSurveyUrl)
	log.Error().Msgf(srv.ListenAndServe().Error())
}
