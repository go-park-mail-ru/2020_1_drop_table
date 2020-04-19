package staff

import (
	"2020_1_drop_table/internal/app/staff/models"
	"context"
)

type Repository interface {
	Add(ctx context.Context, st models.Staff) (models.Staff, error)
	GetByEmail(ctx context.Context, email string) (models.Staff, error)
	GetByID(ctx context.Context, id int) (models.Staff, error)
	Update(ctx context.Context, newStaff models.SafeStaff) (models.SafeStaff, error)
	AddUuid(ctx context.Context, uuid string, id int) error
	CheckIsOwner(ctx context.Context, staffId int) (bool, error)
	DeleteUuid(ctx context.Context, uuid string) error
	GetCafeId(ctx context.Context, uuid string) (int, error)
	GetStaffListByOwnerId(ctx context.Context, ownerId int) (map[string][]models.StaffByOwnerResponse, error)
	DeleteStaff(ctx context.Context, staffId int) error
	UpdatePosition(ctx context.Context, staffId int, newPosition string) error
}
