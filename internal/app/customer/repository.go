package customer

import (
	"2020_1_drop_table/internal/app/customer/models"
	"context"
)

type Repository interface {
	Add(ctx context.Context, cu models.Customer) (models.Customer, error)
	SetLoyaltyPoints(ctx context.Context, points string, customerID string) (models.Customer, error)
	GetByID(ctx context.Context, customerID string) (models.Customer, error)
	DeleteByID(ctx context.Context, customerID string) error
}
