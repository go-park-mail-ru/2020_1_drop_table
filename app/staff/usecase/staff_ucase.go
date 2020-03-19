package usecase

import (
	"2020_1_drop_table/app"
	globalModels "2020_1_drop_table/app/models"
	"2020_1_drop_table/app/staff"
	"2020_1_drop_table/app/staff/models"
	"2020_1_drop_table/projectConfig"
	"2020_1_drop_table/validators"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	uuid "github.com/nu7hatch/gouuid"
	"io"
	"mime/multipart"
	"os"
	"strings"
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

	_, existed, err := s.staffRepo.GetByEmailAndPassword(ctx, newStaff.Email, newStaff.Password)
	if existed {
		return models.SafeStaff{}, globalModels.ErrExisted
	}
	if err != sql.ErrNoRows {
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

	res, err := s.staffRepo.GetById(ctx, id)

	if err != nil {
		return models.SafeStaff{}, err
	}

	return app.GetSafeStaff(res), nil
}

func (s *staffUsecase) Update(c context.Context, newStaff models.SafeStaff) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	session := ctx.Value("session").(*sessions.Session)

	staffID, found := session.Values["userID"]
	if !found || staffID != newStaff.StaffID {
		return globalModels.ErrForbidden
	}

	newStaff.EditedAt = time.Now()

	return s.staffRepo.Update(ctx, newStaff)
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

	staffObj, existed, err := s.staffRepo.GetByEmailAndPassword(ctx, form.Email, form.Password)
	if !existed {
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

func (s *staffUsecase) SaveFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {

	defer file.Close()

	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	uString := u.String()
	folderName := []rune(uString)[:3]
	separatedFilename := strings.Split(header.Filename, ".")
	if len(separatedFilename) <= 1 {
		err := errors.New("bad filename")
		return "", err
	}
	fileType := separatedFilename[len(separatedFilename)-1]

	path := fmt.Sprintf("%s/%s/%s", projectConfig.MediaFolder, folder, string(folderName))
	filename := fmt.Sprintf("%s.%s", uString, fileType)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", nil
	}

	fullFilename := fmt.Sprintf("%s/%s", path, filename)

	f, err := os.OpenFile(fullFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	return fullFilename, err
}
