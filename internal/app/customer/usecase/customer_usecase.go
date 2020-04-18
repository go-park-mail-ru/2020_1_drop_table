package usecase

import (
	"2020_1_drop_table/internal/app/customer"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/staff"
	"context"
	"database/sql"
	"time"
)

type customerUsecase struct {
	customerRepo   customer.Repository
	staffUsecase   staff.Usecase
	contextTimeout time.Duration
}

func NewCustomerUsecase(c customer.Repository, s staff.Usecase, timeout time.Duration) customer.Usecase {
	return &customerUsecase{
		contextTimeout: timeout,
		customerRepo:   c,
		staffUsecase:   s,
	}
}

func (u customerUsecase) GetPoints(ctx context.Context, uuid string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	cust, err := u.customerRepo.GetByID(ctx, uuid)
	if err != nil {
		return -1, err
	}
	return cust.Points, err

}

func (u customerUsecase) SetPoints(ctx context.Context, uuid string, points int) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	requestStaff, err := u.staffUsecase.GetFromSession(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return globalModels.ErrForbidden
		}
		return err
	}
	targetCustomer, err := u.customerRepo.GetByID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return globalModels.ErrNotFound
		}
		return err
	}
	if requestStaff.CafeId != targetCustomer.CafeID {
		return globalModels.ErrForbidden
	}
	if points < 0 {
		return globalModels.ErrPointsError
	}
	_, err = u.customerRepo.SetLoyaltyPoints(ctx, points, uuid)
	return err
}

func calcSalePercent(sum float32, uuid string) float32 {
	//TODO implement
	return 0.5

}

func (u customerUsecase) GetSale(ctx context.Context, sum float32, uuid string) (float32, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	cust, err := u.customerRepo.GetByID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, globalModels.ErrNotFound
		}
		return 0, err
	}

	requestStaff, err := u.staffUsecase.GetFromSession(ctx)
	if err != nil || requestStaff.CafeId != cust.CafeID {
		return 0, globalModels.ErrForbidden
	}

	oldSum := cust.Sum
	salePercent := calcSalePercent(oldSum, uuid)
	sumWithSale := sum * salePercent
	err = u.customerRepo.IncrementSum(ctx, sumWithSale, uuid)
	return sumWithSale, err

}
