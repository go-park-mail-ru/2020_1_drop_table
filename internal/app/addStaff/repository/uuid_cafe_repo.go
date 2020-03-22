package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UuidCafeRepository struct {
	Conn *sqlx.DB
}

func NewUuidCafeRepository(conn *sqlx.DB) UuidCafeRepository {
	Storage := UuidCafeRepository{conn}
	return Storage
}

func (p *UuidCafeRepository) Add(ctx context.Context, uuid string, id int) error {
	query := `INSERT into UuidCafeRepository(uuid, cafeId) VALUES ($1,$2)`
	_, err := p.Conn.ExecContext(ctx, query, uuid, id)
	return err
}
