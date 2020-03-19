package usecase

import (
	"2020_1_drop_table/app/cafe"
	"2020_1_drop_table/app/cafe/models"
	globalModels "2020_1_drop_table/app/models"
	"2020_1_drop_table/validators"
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"time"
)

type cafeUsecase struct {
	cafeRepo       cafe.Repository
	contextTimeout time.Duration
}

func NewCafeUsecase(c cafe.Repository, timeout time.Duration) cafe.Usecase {
	return &cafeUsecase{
		cafeRepo:       c,
		contextTimeout: timeout,
	}
}

//ToDo make permission only for owners
func (cu *cafeUsecase) Add(c context.Context, newCafe models.Cafe) (models.Cafe, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	session := ctx.Value("session").(*sessions.Session)

	staffInterface, found := session.Values["userID"]
	staffID, ok := staffInterface.(int)

	if !found || !ok || staffID <= 0 {
		return models.Cafe{}, globalModels.ErrForbidden
	}

	newCafe.StaffID = staffID

	validation, _, err := validators.GetValidator()
	if err != nil {
		return models.Cafe{}, fmt.Errorf("HttpResponse in validator: %cu", err.Error())
	}

	if err := validation.Struct(newCafe); err != nil {
		return models.Cafe{}, err
	}

	return cu.cafeRepo.Add(ctx, newCafe)
}

func (cu *cafeUsecase) GetByOwnerID(c context.Context) ([]models.Cafe, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	session := ctx.Value("session").(*sessions.Session)

	staffInterface, found := session.Values["userID"]
	staffID, ok := staffInterface.(int)

	if !found || !ok || staffID <= 0 {
		return make([]models.Cafe, 0), globalModels.ErrForbidden
	}

	return cu.cafeRepo.GetByOwnerID(ctx, staffID)
}

func (cu *cafeUsecase) GetByID(c context.Context, id int) (models.Cafe, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	return cu.cafeRepo.GetByID(ctx, id)
}

//ToDo make permission only for owners
func (cu *cafeUsecase) Update(c context.Context, newCafe models.Cafe) error {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	oldCafe, err := cu.cafeRepo.GetByID(ctx, newCafe.CafeID)
	if err != nil {
		return err
	}

	session := ctx.Value("session").(*sessions.Session)

	staffInterface, found := session.Values["userID"]
	staffID, ok := staffInterface.(int)

	if !found || !ok || oldCafe.StaffID != staffID {
		return globalModels.ErrForbidden
	}

	newCafe.StaffID = staffID
	if oldCafe.StaffID != newCafe.StaffID {
		return globalModels.ErrInvalidAction
	}

	return cu.cafeRepo.Update(ctx, newCafe)
}
