package cafe

import (
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
)

type Repository interface {
	Add(ctx context.Context, ca models.Cafe) (models.Cafe, error)
	GetByID(ctx context.Context, id int) (models.Cafe, error)
	GetByOwnerID(ctx context.Context, staffID int) ([]models.Cafe, error)
	Update(ctx context.Context, newCafe models.Cafe) (models.Cafe, error)
	GetAllCafes(ctx context.Context, since int, limit int) ([]models.Cafe, error)
}
