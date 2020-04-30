package staff

import (
	"2020_1_drop_table/internal/microservices/staff/models"
	"context"
)

type StaffClientInterface interface {
	GetFromSession(ctx context.Context) (models.SafeStaff, error)
	GetById(ctx context.Context, id int) (models.SafeStaff, error)
	AddSessionInMetadata(ctx context.Context) context.Context
}
