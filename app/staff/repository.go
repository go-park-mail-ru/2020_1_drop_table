package staff

import (
	"2020_1_drop_table/app/staff/models"
	"context"
)

type Repository interface {
	Add(ctx context.Context, st models.Staff) (models.Staff, error)
	GetByEmailAndPassword(ctx context.Context, email string, password string) (models.Staff, error)
	GetById(ctx context.Context, id int) (models.Staff, error)
	Update(ctx context.Context, newStaff models.SafeStaff) error
}
