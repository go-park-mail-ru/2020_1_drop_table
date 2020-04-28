package main

import (
	"2020_1_drop_table/configs"
	_cafeRepo "2020_1_drop_table/internal/app/cafe/repository"
	grpcServer "2020_1_drop_table/internal/microservices/staff/delivery/grpc/grpc_server"
	_staffRepo "2020_1_drop_table/internal/microservices/staff/repository"
	_staffUsecase "2020_1_drop_table/internal/microservices/staff/usecase"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func main() {
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
	staffRepo := _staffRepo.NewPostgresStaffRepository(conn)
	cafeRepo := _cafeRepo.NewPostgresCafeRepository(conn)
	staffUsecase := _staffUsecase.NewStaffUsecase(&staffRepo, cafeRepo, timeoutContext)
	grpcServer.StartGrpcServer(staffUsecase)
}
