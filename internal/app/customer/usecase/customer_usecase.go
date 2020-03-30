package usecase

import (
	"2020_1_drop_table/internal/app/customer"
	"2020_1_drop_table/internal/app/staff"
	"context"
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
