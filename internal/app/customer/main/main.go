package main

import (
	"2020_1_drop_table/configs"
	_appleRepo "2020_1_drop_table/internal/app/apple_passkit/repository"
	_cafeRepo "2020_1_drop_table/internal/app/cafe/repository"
	server2 "2020_1_drop_table/internal/app/customer/delivery/grpc/server"
	_customerRepo "2020_1_drop_table/internal/app/customer/repository"
	_customerUseCase "2020_1_drop_table/internal/app/customer/usecase"
	_staffRepo "2020_1_drop_table/internal/microservices/staff/repository"
	_staffUsecase "2020_1_drop_table/internal/microservices/staff/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"net"

	"google.golang.org/grpc"

	"fmt"
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
	//ToDo replace with microservices
	staffRepo := _staffRepo.NewPostgresStaffRepository(conn)

	cafeRepo := _cafeRepo.NewPostgresCafeRepository(conn)
	staffUsecase := _staffUsecase.NewStaffUsecase(&staffRepo, cafeRepo, timeoutContext)

	customerRepo := _customerRepo.NewPostgresCustomerRepository(conn)

	applePassKitRepo := _appleRepo.NewPostgresApplePassRepository(conn)

	customerUseCase := _customerUseCase.NewCustomerUsecase(customerRepo, staffUsecase, applePassKitRepo, timeoutContext)

	list, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Err(err)
	}
	server := grpc.NewServer()
	server2.NewCustomerServerGRPC(server, customerUseCase)
	err = server.Serve(list)
	if err != nil {
		log.Error().Msgf(err.Error())
		return
	}
}
