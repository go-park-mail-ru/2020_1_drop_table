package apple_passkit

import (
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"bytes"
	"context"
)

type Usecase interface {
	UpdatePass(c context.Context, pass models.ApplePassDB, cafeID int, publish bool, designOnly bool) error
	GeneratePassObject(c context.Context, cafeID int) (*bytes.Buffer, error)
}
