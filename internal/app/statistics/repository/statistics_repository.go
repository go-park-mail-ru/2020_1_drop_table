package repository

import (
	"2020_1_drop_table/internal/app/statistics"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type postgresStatisticsRepository struct {
	Conn *sqlx.DB
}

func (p postgresStatisticsRepository) GetWorkerDataFromRepo(ctx context.Context, staffId int, limit int, since int) {
	panic("implement me")
}

func (p postgresStatisticsRepository) AddData(jsonData string, time time.Time, clientUUID string, staffID int, cafeID int) error {
	query := `insert into statistics_table (jsonData, time, clientUUID, staffID, cafeID) VALUES ($1,$2,$3,$4,$5)`
	_, err := p.Conn.Exec(query, jsonData, time, clientUUID, staffID, cafeID)
	return err
}

func NewPostgresStatisticsRepository(conn *sqlx.DB) statistics.Repository {
	return &postgresStatisticsRepository{
		Conn: conn,
	}
}
