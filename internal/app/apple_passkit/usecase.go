package apple_passkit

import (
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"bytes"
	"context"
)

type Usecase interface {
	UpdatePass(c context.Context, pass models.ApplePassDB) (models.UpdateResponse, error)

	GeneratePassObject(c context.Context, cafeID int, Type string, published bool) (*bytes.Buffer, error)
	GetPass(c context.Context, cafeID int, Type string, published bool) (map[string]string, error)
	GetImage(c context.Context, imageName string, cafeID int, Type string, published bool) ([]byte, error)
}
