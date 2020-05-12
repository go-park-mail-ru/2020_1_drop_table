package statistics

import (
	"2020_1_drop_table/internal/app/statistics/models"
	"context"
	"time"
)

type Usecase interface {
	AddData(jsonData string, time time.Time, clientUUID string, staffId int, cafeId int) error
	GetWorkerData(ctx context.Context, staffID int, limit int, since int) ([]models.StatisticsStruct, error)
}
