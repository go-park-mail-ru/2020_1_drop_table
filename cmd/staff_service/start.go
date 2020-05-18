package main

import (
	"2020_1_drop_table/configs"
	cafeClient "2020_1_drop_table/internal/app/cafe/delivery/grpc/client"
	"2020_1_drop_table/internal/app/middleware"
	grpcServer "2020_1_drop_table/internal/microservices/staff/delivery/grpc/grpc_server"
	staffHandler "2020_1_drop_table/internal/microservices/staff/delivery/http"
	_staffRepo "2020_1_drop_table/internal/microservices/staff/repository"
	_staffUsecase "2020_1_drop_table/internal/microservices/staff/usecase"
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

	grpcCafeConn, err := grpc.Dial(configs.GRPCCafeUrl, grpc.WithInsecure())
	grpcCafeClient := cafeClient.NewCafeClient(grpcCafeConn)

	staffRepo := _staffRepo.NewPostgresStaffRepository(conn)
	staffUsecase := _staffUsecase.NewStaffUsecase(&staffRepo, grpcCafeClient, timeoutContext)

	go grpcServer.StartStaffGrpcServer(staffUsecase, configs.GRPCStaffUrl)
	staffHandler.NewStaffHandler(r, staffUsecase)

	//static server
	r.PathPrefix(fmt.Sprintf("/%s/", configs.MediaFolder)).Handler(
		http.StripPrefix(fmt.Sprintf("/%s/", configs.MediaFolder),
			http.FileServer(http.Dir(configs.MediaFolder))))

	middleware.AddOptionsRequest(r)

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         configs.HTTPStaffUrl,
		WriteTimeout: configs.Timeouts.WriteTimeout,
		ReadTimeout:  configs.Timeouts.ReadTimeout,
	}
	fmt.Println("staff server started at ", configs.HTTPStaffUrl)
	log.Error().Msgf(srv.ListenAndServe().Error())
}
