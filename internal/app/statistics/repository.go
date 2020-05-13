package statistics

import (
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	"2020_1_drop_table/internal/app/statistics/models"
	"context"
	"time"
)

type Repository interface {
	AddData(jsonData string, time time.Time, clientUUID string, staffID int, cafeId int) error
	GetWorkerDataFromRepo(ctx context.Context, staffId int, limit int, since int) ([]models.StatisticsStruct, error)
	GetGraphsDataFromRepo(ctx context.Context, cafeList []cafeModels.Cafe, typ string, since string, to string)
}
