package usecase

import (
	"2020_1_drop_table/internal/app"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/staff"
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/pkg/qr"
	"2020_1_drop_table/internal/pkg/validators"
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	uuid "github.com/nu7hatch/gouuid"
	"time"
)

type staffUsecase struct {
	staffRepo      staff.Repository
	contextTimeout time.Duration
}

func NewStaffUsecase(s staff.Repository, timeout time.Duration) staff.Usecase {
	return &staffUsecase{
		staffRepo:      s,
		contextTimeout: timeout,
	}
}

func (s *staffUsecase) Add(c context.Context, newStaff models.Staff) (models.SafeStaff, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	newStaff.EditedAt = time.Now()

	validation, _, err := validators.GetValidator()
	if err != nil {
		return models.SafeStaff{}, fmt.Errorf("HttpResponse in validator: %s", err.Error())
	}

	if err := validation.Struct(newStaff); err != nil {
		return models.SafeStaff{}, err
	}

	newStaff.Password = getMD5Hash(newStaff.Password)

	_, err = s.staffRepo.GetByEmailAndPassword(ctx, newStaff.Email, newStaff.Password)
	if err != sql.ErrNoRows {
		return models.SafeStaff{}, globalModels.ErrExisted
	}
	if err != nil && err != sql.ErrNoRows {
		return models.SafeStaff{}, err
	}

	newStaff, err = s.staffRepo.Add(ctx, newStaff)
	if err != nil {
		return models.SafeStaff{}, err
	}

	return app.GetSafeStaff(newStaff), nil
}

func (s *staffUsecase) GetByID(c context.Context, id int) (models.SafeStaff, error) {

	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	session := ctx.Value("session").(*sessions.Session)

	staffID, found := session.Values["userID"]
	if !found || staffID != id {
		return models.SafeStaff{}, globalModels.ErrForbidden
	}

	res, err := s.staffRepo.GetByID(ctx, id)

	if err != nil {
		return models.SafeStaff{}, err
	}

	return app.GetSafeStaff(res), nil
}

func (s *staffUsecase) Update(c context.Context, newStaff models.SafeStaff) (models.SafeStaff, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	session := ctx.Value("session").(*sessions.Session)

	staffID, found := session.Values["userID"]
	if !found || staffID != newStaff.StaffID {
		return models.SafeStaff{}, globalModels.ErrForbidden
	}
	newStaff.EditedAt = time.Now()

	return newStaff, s.staffRepo.Update(ctx, newStaff)
}

func (s *staffUsecase) GetByEmailAndPassword(c context.Context, form models.LoginForm) (models.SafeStaff, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	validation, _, err := validators.GetValidator()
	if err != nil {
		return models.SafeStaff{}, fmt.Errorf("HttpResponse in validator: %s", err.Error())
	}
	if err := validation.Struct(form); err != nil {
		return models.SafeStaff{}, err
	}

	form.Password = getMD5Hash(form.Password)

	staffObj, err := s.staffRepo.GetByEmailAndPassword(ctx, form.Email, form.Password)
	if err == sql.ErrNoRows {
		return models.SafeStaff{}, globalModels.ErrNotFound
	}
	if err != nil {
		return models.SafeStaff{}, globalModels.ErrNotFound
	}

	return app.GetSafeStaff(staffObj), nil
}

func (s *staffUsecase) GetFromSession(c context.Context) (models.SafeStaff, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	session := ctx.Value("session").(*sessions.Session)

	staffID, found := session.Values["userID"]
	if !found || staffID == -1 {
		return models.SafeStaff{}, globalModels.ErrForbidden
	}

	return s.GetByID(ctx, staffID.(int))
}

func (s *staffUsecase) GetQrForStaff(ctx context.Context, idCafe int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	uString := u.String()

	err = s.staffRepo.AddUuid(ctx, uString, idCafe)
	path, err := GenerateQrCode(uString)
	if err != nil {
		return "", err
	}
	return path, nil

}

func GenerateQrCode(uString string) (string, error) {
	link := fmt.Sprintf("/api/v1/staff/addStaff?uuid=%s", uString)
	pathToQr, err := qr.GenerateToFile(link, uString)
	if err != nil {
		return "", err
	}
	return pathToQr, err
}
