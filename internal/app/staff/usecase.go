package staff

import (
	"2020_1_drop_table/internal/app/staff/models"
	"context"
)

type Usecase interface {
	Add(c context.Context, newStaff models.Staff) (models.SafeStaff, error)
	GetByID(c context.Context, id int) (models.SafeStaff, error)
	Update(c context.Context, newStaff models.SafeStaff) (models.SafeStaff, error)
	GetByEmailAndPassword(c context.Context, form models.LoginForm) (models.SafeStaff, error)
	GetFromSession(c context.Context) (models.SafeStaff, error)
	GetQrForStaff(ctx context.Context, idCafe int) (string, error)
	IsOwner(c context.Context, staffId int) (bool, error)
	DeleteQrCodes(uString string) error
	GetCafeId(c context.Context, uuid string) (int, error)
	GetStaffId(c context.Context) (int, error)
	GetStaffListByOwnerId(ctx context.Context, ownerId int) (map[string][]models.StaffByOwnerResponse, error)
	DeleteStaffById(ctx context.Context, staffId int) error
	CheckIfStaffInOwnerCafes(ctx context.Context, requestUser models.SafeStaff, staffId int) (bool, error)
}
