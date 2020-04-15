package main

import (
	"2020_1_drop_table/configs"
	_appleHttpDeliver "2020_1_drop_table/internal/app/apple_passkit/delivery/http"
	_appleRepo "2020_1_drop_table/internal/app/apple_passkit/repository"
	_appleUsecase "2020_1_drop_table/internal/app/apple_passkit/usecase"
	_cafeHttpDeliver "2020_1_drop_table/internal/app/cafe/delivery/http"
	_cafeRepo "2020_1_drop_table/internal/app/cafe/repository"
	_cafeUsecase "2020_1_drop_table/internal/app/cafe/usecase"
	_customerHttpDeliver "2020_1_drop_table/internal/app/customer/delivery/http"
	_customerRepo "2020_1_drop_table/internal/app/customer/repository"
	_customerUseCase "2020_1_drop_table/internal/app/customer/usecase"
	"2020_1_drop_table/internal/app/middleware"
	_staffHttpDeliver "2020_1_drop_table/internal/app/staff/delivery/http"
	_staffRepo "2020_1_drop_table/internal/app/staff/repository"
	_staffUsecase "2020_1_drop_table/internal/app/staff/usecase"
	"2020_1_drop_table/internal/pkg/apple_pass_generator"
	"2020_1_drop_table/internal/pkg/apple_pass_generator/meta"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	redisStore "gopkg.in/boj/redistore.v1"
	"net/http"
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
	_staffHttpDeliver.NewStaffHandler(r, staffUsecase)

	cafeUsecase := _cafeUsecase.NewCafeUsecase(cafeRepo, staffUsecase, timeoutContext)
	_cafeHttpDeliver.NewCafeHandler(r, cafeUsecase)

	applePassGenerator := apple_pass_generator.NewGenerator(
		configs.AppleWWDR, configs.AppleCertificate, configs.AppleKey, configs.ApplePassword)

	customerRepo := _customerRepo.NewPostgresCustomerRepository(conn)

	applePassKitRepo := _appleRepo.NewPostgresApplePassRepository(conn)

	applePassKitUcase := _appleUsecase.NewApplePassKitUsecase(applePassKitRepo, cafeRepo, customerRepo,
		&applePassGenerator, timeoutContext, &meta.Meta{})

	_appleHttpDeliver.NewPassKitHandler(r, applePassKitUcase)

	customerUseCase := _customerUseCase.NewCustomerUsecase(customerRepo, staffUsecase, timeoutContext)
	_customerHttpDeliver.NewCustomerHandler(r, customerUseCase)

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
}
