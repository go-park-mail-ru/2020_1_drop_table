package usecase

import (
	cafeClient "2020_1_drop_table/internal/app/cafe/delivery/grpc/client"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/statistics"
	"2020_1_drop_table/internal/app/statistics/models"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	"context"
	"errors"
	"fmt"
	"time"
)

type statisticsUsecase struct {
	statisticsRepo statistics.Repository
	staffClient    staffClient.StaffClientInterface
	cafeClient     cafeClient.CafeGRPCClientInterface
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

func NewStatisticsUsecase(st statistics.Repository, staffClient staffClient.StaffClientInterface, cafeClient cafeClient.CafeGRPCClientInterface, timeout time.Duration) *statisticsUsecase {
	return &statisticsUsecase{
		statisticsRepo: st,
		contextTimeout: timeout,
		staffClient:    staffClient,
		cafeClient:     cafeClient,
	}
}

func (s statisticsUsecase) GetDataForGraphs(ctx context.Context, typ string, since string, to string) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	requestUser, err := s.staffClient.GetFromSession(ctx)
	if err != nil {
		return globalModels.ErrForbidden
	}
	cafes, err := s.cafeClient.GetByOwnerId(ctx, requestUser.StaffID)
	rawData, err := s.statisticsRepo.GetGraphsDataFromRepo(ctx, cafes, typ, since, to)
	fmt.Println(rawData)
	jsonify(rawData)
	return nil
}

func jsonify(data []models.StatisticsGraphRawStruct) {
	var m = make(map[int]interface{})
	var prevElem = data[0]
	for _, el := range data {
		fmt.Println(el)
		if el.CafeId > prevElem.CafeId {
			m[el.CafeId] = el

		}
		prevElem = el
	}
	fmt.Println("map: ", m)
}
