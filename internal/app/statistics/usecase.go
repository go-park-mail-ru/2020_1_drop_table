package statistics

import (
	"context"
	"time"
)

type Usecase interface {
	AddData(jsonData string, time time.Time, clientUUID string, staffID int, cafeID int) error
	GetWorkerData(ctx context.Context, staffID int, limit int, since int)
}
