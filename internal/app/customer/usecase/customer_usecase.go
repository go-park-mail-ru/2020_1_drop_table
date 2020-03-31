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
			return globalModels.RightsError
		}
		return err
	}
	targetCustomer, err := u.customerRepo.GetByID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return globalModels.BadUuid
		}
		return err
	}
	if requestStaff.CafeId != targetCustomer.CafeID {
		return globalModels.RightsError
	}
	if points < 0 {
		return globalModels.PointsError
	}
	_, err = u.customerRepo.SetLoyaltyPoints(ctx, points, uuid)
	return err
}
