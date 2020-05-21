package client

import (
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
)

type CafeGRPCClientInterface interface {
	GetByID(ctx context.Context, id int) (models.Cafe, error)
	GetByOwnerId(ctx context.Context, ownerID int) ([]models.Cafe, error)
}
