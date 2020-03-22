package usecase

import (
	"2020_1_drop_table/internal/app/addStaff"
	"2020_1_drop_table/internal/pkg/qr"
	"context"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"time"
)

type AddStaffUsecase struct {
	uuidCafeRepository addStaff.Repository
	contextTimeout     time.Duration
}

func newAddStaffUsecase(s addStaff.Repository, timeout time.Duration) addStaff.Usecase {
	return &AddStaffUsecase{
		uuidCafeRepository: s,
		contextTimeout:     timeout,
	}
}

func (s *AddStaffUsecase) GetQrForStaff(ctx context.Context, idCafe int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	uString := u.String()

	err = s.uuidCafeRepository.Add(ctx, uString, idCafe)
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
