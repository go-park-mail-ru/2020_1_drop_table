package customer

import (
	"context"
)

type Usecase interface {
	GetPoints(ctx context.Context, uuid string) (int, error)
	SetPoints(ctx context.Context, uuid string, points int) error
	GetSale(ctx context.Context, sum float32, uuid string) (float32, error)
}
