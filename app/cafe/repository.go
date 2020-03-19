package cafe

import (
	"2020_1_drop_table/app/cafe/models"
	"context"
)

type Repository interface {
	Add(ctx context.Context, ca models.Cafe) (models.Cafe, error)
	GetByID(ctx context.Context, id int) (models.Cafe, error)
	GetByOwnerID(ctx context.Context, staffID int) ([]models.Cafe, error)
	Update(ctx context.Context, newCafe models.Cafe) error
}
