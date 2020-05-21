package cafe

import (
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
)

type Usecase interface {
	Add(c context.Context, newCafe models.Cafe) (models.Cafe, error)
	GetByOwnerID(c context.Context) ([]models.Cafe, error)
	GetByOwnerIDWithOwnerID(ctx context.Context, ownerID int) ([]models.Cafe, error)
	GetByID(c context.Context, id int) (models.Cafe, error)
	Update(c context.Context, newCafe models.Cafe) (models.Cafe, error)
	GetAllCafes(ctx context.Context, since int, limit int, search string) ([]models.Cafe, error)
	GetCafeSortedByRadius(ctx context.Context, latitude string, longitude string, radius string) ([]models.Cafe, error)
}
