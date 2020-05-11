package usecase

import (
	"2020_1_drop_table/internal/app/statistics"
	"context"
	"time"
)

type statisticsUsecase struct {
	statisticsRepo statistics.Repository
	contextTimeout time.Duration
}

func (s statisticsUsecase) AddData(jsonData string, time time.Time, clientUUID string, staffID int, cafeID int) error {
	return s.statisticsRepo.AddData(jsonData, time, clientUUID, staffID, cafeID)

}

func (s statisticsUsecase) GetWorkerData(ctx context.Context, staffID int, limit int, since int) {

}

func NewStatisticsUsecase(st statistics.Repository, timeout time.Duration) *statisticsUsecase {
	return &statisticsUsecase{
		statisticsRepo: st,
		contextTimeout: timeout,
	}
}
