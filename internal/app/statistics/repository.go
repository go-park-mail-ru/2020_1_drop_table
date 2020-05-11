package statistics

import (
	"context"
	"time"
)

type Repository interface {
	AddData(jsonData string, time time.Time, clientUUID string, staffID int, cafeID int) error
	GetWorkerDataFromRepo(ctx context.Context, staffId int, limit int, since int)
}
