package usecase

import (
	"2020_1_drop_table/internal/app/cafe"
	"2020_1_drop_table/internal/app/cafe/models"
	globalModels "2020_1_drop_table/internal/app/models"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	geo "2020_1_drop_table/internal/pkg/google_geocoder"
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type cafeUsecase struct {
	cafeRepo        cafe.Repository
	staffGrpcClient staffClient.StaffClientInterface
	contextTimeout  time.Duration
	geoCoder        geo.GoogleGeoCoder
}

func (cu *cafeUsecase) GetCafeSortedByRadius(ctx context.Context, latitude string, longitude string, radius string) ([]models.CafeWithGeoData, error) {

	return cu.cafeRepo.GetCafeSortedByRadius(ctx, latitude, longitude, radius)
}

func (cu *cafeUsecase) GetByOwnerIDWithOwnerID(ctx context.Context, ownerID int) ([]models.Cafe, error) {
	return cu.cafeRepo.GetByOwnerID(ctx, ownerID)
}

func NewCafeUsecase(c cafe.Repository, stClient staffClient.StaffClientInterface,
	timeout time.Duration, geoCoder geo.GoogleGeoCoder) cafe.Usecase {
	return &cafeUsecase{
		cafeRepo:        c,
		contextTimeout:  timeout,
		staffGrpcClient: stClient,
		geoCoder:        geoCoder,
	}
}

func (cu *cafeUsecase) checkIsOwnerById(c context.Context, staffID int) (bool, error) {
	staffObj, err := cu.staffGrpcClient.GetById(c, staffID)

	if err != nil {
		return false, err
	}

	return staffObj.IsOwner, nil
}

func (cu *cafeUsecase) Add(c context.Context, newCafe models.Cafe) (models.Cafe, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	session := ctx.Value("session").(*sessions.Session)

	staffInterface, found := session.Values["userID"]
	staffID, ok := staffInterface.(int)

	if !found || !ok || staffID <= 0 {
		return models.Cafe{}, globalModels.ErrForbidden
	}

	isOwner, err := cu.checkIsOwnerById(c, staffID)
	if err != nil {
		return models.Cafe{}, err
	}
	if !isOwner {
		return models.Cafe{}, globalModels.ErrForbidden
	}

	newCafe.StaffID = staffID

	validation := validator.New()

	if err := validation.Struct(newCafe); err != nil {
		return models.Cafe{}, err
	}

	if newCafe.Address != "" {
		geoInfo, err := cu.geoCoder.GetGeoByAddress(newCafe.Address)
		fmt.Println(geoInfo, err)
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

func (cu *cafeUsecase) Update(c context.Context, newCafe models.Cafe) (models.Cafe, error) {
	ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

	oldCafe, err := cu.cafeRepo.GetByID(ctx, newCafe.CafeID)
	if err != nil {
		return models.Cafe{}, err
	}

	session := ctx.Value("session").(*sessions.Session)

	staffInterface, found := session.Values["userID"]
	staffID, ok := staffInterface.(int)

	if !found || !ok || oldCafe.StaffID != staffID {
		return models.Cafe{}, globalModels.ErrForbidden
	}

	newCafe.StaffID = staffID
	if oldCafe.StaffID != newCafe.StaffID {
		return models.Cafe{}, globalModels.ErrInvalidAction
	}

	validation := validator.New()

	if err := validation.Struct(newCafe); err != nil {
		return models.Cafe{}, err
	}

	return cu.cafeRepo.Update(ctx, newCafe)
}

func (cu *cafeUsecase) GetAllCafes(ctx context.Context, since int, limit int) ([]models.Cafe, error) {
	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	cafes, err := cu.cafeRepo.GetAllCafes(ctx, since, limit)
	return cafes, err
}
