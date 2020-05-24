package usecase

import (
	"2020_1_drop_table/internal/app/apple_passkit"
	"2020_1_drop_table/internal/app/customer"
	"2020_1_drop_table/internal/app/customer/models"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/statistics"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	loyaltySystems "2020_1_drop_table/internal/pkg/apple_pass_generator/loyalty_systems"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type customerUsecase struct {
	customerRepo      customer.Repository
	passKitRepo       apple_passkit.Repository
	staffClient       staffClient.StaffClientInterface
	statisticsUsecase statistics.Usecase
	contextTimeout    time.Duration
}

func NewCustomerUsecase(c customer.Repository, p apple_passkit.Repository, clientInterface staffClient.StaffClientInterface,
	timeout time.Duration, statUsecase statistics.Usecase) customer.Usecase {

	return &customerUsecase{
		contextTimeout:    timeout,
		passKitRepo:       p,
		customerRepo:      c,
		staffClient:       clientInterface,
		statisticsUsecase: statUsecase,
	}
}

func (u customerUsecase) GetCustomer(ctx context.Context, uuid string) (models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	//ToDo make permission only for staff after adding statistics
	targetCustomer, err := u.customerRepo.GetByID(ctx, uuid)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Customer{}, globalModels.ErrNotFound
		}
		return models.Customer{}, err
	}

	return targetCustomer, err
}

func (u customerUsecase) GetPoints(ctx context.Context, uuid string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	cust, err := u.customerRepo.GetByID(ctx, uuid)
	if err != nil {
		return "", err
	}
	return cust.Points, err
}

func (u customerUsecase) SetPoints(ctx context.Context, uuid string, points string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	currTime := time.Now()

	requestStaff, err := u.staffClient.GetFromSession(ctx)
	if err != nil {
		fmt.Println("err here: ", err)
		if err == sql.ErrNoRows {
			return globalModels.ErrForbidden
		}
		return err
	}

	targetCustomer, err := u.customerRepo.GetByID(ctx, uuid)
	if err != nil {
		fmt.Println("err here 2: ", err)
		if err == sql.ErrNoRows {
			return globalModels.ErrNotFound
		}
		return err
	}
	if requestStaff.CafeId != targetCustomer.CafeID {
		return globalModels.ErrForbidden
	}

	pass, err := u.passKitRepo.GetPassByCafeID(ctx, targetCustomer.CafeID, targetCustomer.Type, true)
	if err != nil {
		return err
	}

	loyaltySystem, ok := loyaltySystems.LoyaltySystems[targetCustomer.Type]
	if !ok {
		return err
	}

	newPoints, err := loyaltySystem.SettingPoints(pass.LoyaltyInfo, targetCustomer.Points, points)
	if err != nil {
		return err
	}

	//todo check if this work all together
	err = u.statisticsUsecase.AddData(newPoints, currTime, targetCustomer.CustomerID, requestStaff.StaffID, requestStaff.CafeId)
	if err != nil {

		return err
	}

	_, err = u.customerRepo.SetLoyaltyPoints(ctx, newPoints, uuid)
	return err
}

func (u customerUsecase) Add(ctx context.Context, newCustomer models.Customer) (models.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.customerRepo.Add(ctx, newCustomer)
}
