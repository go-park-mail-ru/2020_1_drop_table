package usecase

import (
	"2020_1_drop_table/internal/app/statistics"
	"2020_1_drop_table/internal/app/statistics/models"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	"context"
	"errors"
	"time"
)

type statisticsUsecase struct {
	statisticsRepo statistics.Repository
	staffClient    staffClient.StaffClientInterface
	contextTimeout time.Duration
}

func (s statisticsUsecase) AddData(jsonData string, time time.Time, clientUUID string, staffId int, cafeId int) error {
	return s.statisticsRepo.AddData(jsonData, time, clientUUID, staffId, cafeId)

}

func (s statisticsUsecase) GetWorkerData(ctx context.Context, staffID int, limit int, since int) ([]models.StatisticsStruct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	requestStaff, err := s.staffClient.GetFromSession(ctx)

	//todo permissions check if staffId working on request user
	if err != nil || !requestStaff.IsOwner {
		message := "Request user not owner or not authorized"
		return nil, errors.New(message)
	}

	data, err := s.statisticsRepo.GetWorkerDataFromRepo(ctx, staffID, limit, since)

	return data, err
}

func NewStatisticsUsecase(st statistics.Repository, staffClient staffClient.StaffClientInterface, timeout time.Duration) *statisticsUsecase {
	return &statisticsUsecase{
		statisticsRepo: st,
		contextTimeout: timeout,
		staffClient:    staffClient,
	}
}
