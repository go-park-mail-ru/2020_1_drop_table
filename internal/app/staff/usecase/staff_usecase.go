package usecase

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app"
	"2020_1_drop_table/internal/app/cafe"
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/staff"
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/pkg/qr"
	"2020_1_drop_table/internal/pkg/validators"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type staffUsecase struct {
	staffRepo      staff.Repository
	cafeRepo       cafe.Repository
	contextTimeout time.Duration
}

func NewStaffUsecase(s staff.Repository, cafeRepo cafe.Repository, timeout time.Duration) staff.Usecase {
	return &staffUsecase{
		staffRepo:      s,
		contextTimeout: timeout,
		cafeRepo:       cafeRepo,
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
	staffId := s.GetStaffId(ctx)
	owner, err := s.GetByID(ctx, staffId)
	if err != nil {
		message := fmt.Sprintf("Cant find Staff in SessionStorage because of -> %s", err)
		log.Error().Msgf(message)
		return "", errors.New(message)
	}

	ownerCafe, err := s.cafeRepo.GetByOwnerID(ctx, owner.StaffID)
	if err != nil {
		message := fmt.Sprintf("Cant find cafe with this owner because of -> %s", err)
		if err == sql.ErrNoRows {
			message = fmt.Sprintf("Cant find cafe with this owner")
		}
		log.Error().Msgf(message)
		return "", errors.New(message)
	}

	isIn := isCafeInCafeList(idCafe, ownerCafe)
	if owner.IsOwner && isIn {

		u, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		uString := u.String()

		err = s.staffRepo.AddUuid(ctx, uString, idCafe)
		path, err := generateQRCode(uString)
		if err != nil {
			return "", err
		}
		return path, nil
	}
	message := fmt.Sprintf("User is not owner of this cafe")
	log.Error().Msgf(message)
	return "", errors.New(message)
}

func isCafeInCafeList(idCafe int, ownersCafe []cafeModels.Cafe) bool {
	for _, cafe := range ownersCafe {
		if cafe.StaffID == idCafe {
			fmt.Println(cafe)
			return true
		}
	}
	return false
}

func (s *staffUsecase) DeleteQrCodes(uString string) error {
	pathToQr := configs.MediaFolder + "/qr/" + uString + ".png"
	os.Remove(pathToQr)
	err := s.staffRepo.DeleteUuid(context.TODO(), uString)
	return err

}

func generateQRCode(uString string) (string, error) {
	link := fmt.Sprintf("%s/addStaff?uuid=%s", configs.FrontEndUrl, uString)
	pathToQr, err := qr.GenerateToFile(link, uString)
	if err != nil {
		return "", err
	}
	return pathToQr, err
}

func (s *staffUsecase) IsOwner(c context.Context, staffId int) (bool, error) {
	return s.staffRepo.CheckIsOwner(c, staffId)
}

func (s *staffUsecase) GetCafeId(c context.Context, uuid string) (int, error) {
	return s.staffRepo.GetCafeId(c, uuid)
}

func (s *staffUsecase) GetStaffId(c context.Context) int { //TODO взять Димин код
	session := c.Value("session").(*sessions.Session)
	staffID, _ := session.Values["userID"]
	return staffID.(int)

}
