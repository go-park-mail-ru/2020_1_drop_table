package main

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app/middleware"
	grpcServer "2020_1_drop_table/internal/microservices/staff/delivery/grpc/grpc_server"
	_staffRepo "2020_1_drop_table/internal/microservices/staff/repository"
	_staffUsecase "2020_1_drop_table/internal/microservices/staff/usecase"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	redisStore "gopkg.in/boj/redistore.v1"
)

func main() {
	r := mux.NewRouter()
	timeoutContext := configs.Timeouts.ContextTimeout
	//Middleware
	var CookieStore, err = redisStore.NewRediStore(
		configs.RedisPreferences.Size,
		configs.RedisPreferences.Network,
		configs.RedisPreferences.Address,
		configs.RedisPreferences.Password,
		configs.RedisPreferences.SecretKey)

	middleware.NewMiddleware(r, CookieStore)

	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s",
		configs.PostgresPreferences.User,
		configs.PostgresPreferences.Password,
		configs.PostgresPreferences.Port)

	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Error().Msgf(err.Error())
		return
	}
	staffRepo := _staffRepo.NewPostgresStaffRepository(conn)
	staffUsecase := _staffUsecase.NewStaffUsecase(&staffRepo, timeoutContext)
	grpcServer.StartGrpcServer(staffUsecase)
}
