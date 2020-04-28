package customer

import (
	"2020_1_drop_table/internal/app/customer/models"
	"context"
)

type Usecase interface {
	GetPoints(ctx context.Context, uuid string) (string, error)
	SetPoints(ctx context.Context, uuid, points string) error
	GetCustomer(ctx context.Context, uuid string) (models.Customer, error)
	Add(ctx context.Context, newCustomer models.Customer) (models.Customer, error)
}
