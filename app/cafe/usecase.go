package cafe

import (
	"2020_1_drop_table/app/cafe/models"
	"context"
)

type Usecase interface {
	Add(c context.Context, newCafe models.Cafe) (models.Cafe, error)
	GetByOwnerID(c context.Context) ([]models.Cafe, error)
	GetByID(c context.Context, id int) (models.Cafe, error)
	Update(c context.Context, newCafe models.Cafe) error
}
