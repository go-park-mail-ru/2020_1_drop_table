package usecase

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app"
	cafeClient "2020_1_drop_table/internal/app/cafe/delivery/grpc/client"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/microservices/staff"
	"2020_1_drop_table/internal/microservices/staff/models"
	"2020_1_drop_table/internal/pkg/hasher"
	"2020_1_drop_table/internal/pkg/qr"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"time"
)

type staffUsecase struct {
	staffRepo      staff.Repository
	cafeClient     cafeClient.CafeGRPCClientInterface
	contextTimeout time.Duration
}

func NewStaffUsecase(s staff.Repository, cafeClient cafeClient.CafeGRPCClientInterface, timeout time.Duration) staff.Usecase {
	return &staffUsecase{
		staffRepo:      s,
		contextTimeout: timeout,
		cafeClient:     cafeClient,
	}
}

func (s *staffUsecase) Add(c context.Context, newStaff models.Staff) (models.SafeStaff, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	newStaff.EditedAt = time.Now()

	validation := validator.New()

	if err := validation.Struct(newStaff); err != nil {
		return models.SafeStaff{}, err
	}

	_, err := s.staffRepo.GetByEmail(ctx, newStaff.Email)

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

	validation := validator.New()
	if err := validation.Struct(newStaff); err != nil {
		return models.SafeStaff{}, err
	}

	return s.staffRepo.Update(ctx, newStaff)
}

func (s *staffUsecase) GetByEmailAndPassword(c context.Context, form models.LoginForm) (models.SafeStaff, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	validation := validator.New()
	if err := validation.Struct(form); err != nil {
		return models.SafeStaff{}, err
	}

	staffObj, err := s.staffRepo.GetByEmail(ctx, form.Email)

	if err == sql.ErrNoRows {
		return models.SafeStaff{}, globalModels.ErrNotFound
	}
	if err != nil {
		return models.SafeStaff{}, globalModels.ErrNotFound
	}

	if !hasher.CheckWithHash(staffObj.Password, form.Password) {
		return models.SafeStaff{}, globalModels.ErrUnauthorized
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

func (s *staffUsecase) GetQrForStaff(ctx context.Context, idCafe int, position string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	staffId, err := s.GetStaffId(ctx)
	if err != nil {
		message := fmt.Sprintf("Cant find stuff in GET Params of -> %s", err)
		log.Error().Msgf(message)
		return "", errors.New(message)
	}
	owner, err := s.GetByID(ctx, staffId)
	if err != nil {
		message := fmt.Sprintf("Cant find Staff in SessionStorage because of -> %s", err)
		log.Error().Msgf(message)
		return "", errors.New(message)
	}
	ownerCafe, err := s.cafeClient.GetByID(ctx, idCafe)
	if err != nil {
		message := fmt.Sprintf("Cant find cafe with this owner because of -> %s", err)
		if err == sql.ErrNoRows {
			message = fmt.Sprintf("User is not owner of cafe")
		}
		log.Error().Msgf(message)
		return "", errors.New(message)
	}
	if owner.IsOwner && ownerCafe.StaffID == owner.StaffID {

		u, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		uString := u.String()
		err = s.staffRepo.AddUuid(ctx, uString, idCafe)
		path, err := generateQRCode(uString, position)
		if err != nil {
			return "", err
		}
		return path, nil
	}
	message := fmt.Sprintf("User is not owner of this cafe")
	log.Error().Msgf(message)
	return "", errors.New(message)
}

func (s *staffUsecase) DeleteQrCodes(uString string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), s.contextTimeout)
	defer cancel()
	pathToQr := configs.MediaFolder + "/qr/" + uString + ".png"
	err := os.Remove(pathToQr)
	if err != nil {
		return err
	}
	err = s.staffRepo.DeleteUuid(ctx, uString)
	return err

}

func generateQRCode(uString string, position string) (string, error) {
	link := fmt.Sprintf("%s/addStaff?uuid=%s&position=%s", configs.FrontEndUrl, uString, position)
	pathToQr, err := qr.GenerateToFile(link, uString)
	pathToQr = configs.HTTPStaffUrl + "/" + pathToQr
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

func (s *staffUsecase) GetStaffId(c context.Context) (int, error) {
	session := c.Value("session").(*sessions.Session)

	staffID, ok := session.Values["userID"]
	if !ok {
		return -1, errors.New("no userID in session")
	}

	id, ok := staffID.(int)
	if !ok {
		return -1, errors.New("userID is not int")
	}
	return id, nil

}

func (s *staffUsecase) GetStaffListByOwnerId(ctx context.Context, ownerId int) (map[string][]models.StaffByOwnerResponse, error) {
	requestUser, err := s.GetFromSession(ctx)
	if err != nil {
		emptMap := make(map[string][]models.StaffByOwnerResponse)
		return emptMap, err
	}
	if requestUser.IsOwner && requestUser.StaffID == ownerId {
		return s.staffRepo.GetStaffListByOwnerId(ctx, ownerId)
	}
	emptMap := make(map[string][]models.StaffByOwnerResponse)

	return emptMap, globalModels.ErrForbidden
}

func (s *staffUsecase) CheckIfStaffInOwnerCafes(ctx context.Context, requestUser models.SafeStaff, staffId int) (bool, error) {
	if requestUser.StaffID == staffId && requestUser.IsOwner == true {
		return true, nil
	}
	staffList, err := s.GetStaffListByOwnerId(ctx, requestUser.StaffID)
	if err != nil {
		return false, globalModels.ErrForbidden
	}
	for _, cafeStaffList := range staffList {
		for _, staf := range cafeStaffList {
			if *staf.StaffId == staffId {
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *staffUsecase) DeleteStaffById(ctx context.Context, staffId int) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	requestUser, err := s.GetFromSession(ctx)
	if err != nil {
		return globalModels.ErrForbidden
	}
	isIn, err := s.CheckIfStaffInOwnerCafes(ctx, requestUser, staffId)
	if err != nil || !isIn {
		return globalModels.ErrForbidden
	}
	err = s.staffRepo.DeleteStaff(ctx, staffId)
	return err
}

func (s *staffUsecase) UpdatePosition(ctx context.Context, staffId int, newPosition string) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	requestUser, err := s.GetFromSession(ctx)
	if err != nil {
		return globalModels.ErrForbidden
	}
	isIn, err := s.CheckIfStaffInOwnerCafes(ctx, requestUser, staffId)
	if err != nil || !isIn {
		return globalModels.ErrForbidden
	}
	err = s.staffRepo.UpdatePosition(ctx, staffId, newPosition)
	return err
}
