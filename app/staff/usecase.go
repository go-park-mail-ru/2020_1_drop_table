package staff

import (
	"2020_1_drop_table/app/staff/models"
	"context"
	"mime/multipart"
)

type Usecase interface {
	Add(c context.Context, newStaff models.Staff) (models.SafeStaff, error)
	GetByID(c context.Context, id int) (models.SafeStaff, error)
	Update(c context.Context, newStaff models.SafeStaff) error
	GetByEmailAndPassword(c context.Context, form models.LoginForm) (models.SafeStaff, error)
	SaveFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error)
	GetFromSession(c context.Context) (models.SafeStaff, error)
}
