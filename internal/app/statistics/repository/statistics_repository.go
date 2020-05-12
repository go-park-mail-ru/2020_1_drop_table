package repository

import (
	"2020_1_drop_table/internal/app/statistics"
	"2020_1_drop_table/internal/app/statistics/models"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type postgresStatisticsRepository struct {
	Conn *sqlx.DB
}

func (p postgresStatisticsRepository) GetWorkerDataFromRepo(ctx context.Context, staffId int, limit int, since int) ([]models.StatisticsStruct, error) {
	query := `SELECT * from statistics_table where staffID=$1 order by time LIMIT $2 OFFSET $3 `
	var res []models.StatisticsStruct

	err := p.Conn.SelectContext(ctx, &res, query, staffId, limit, since)

	return res, err

}

func (p postgresStatisticsRepository) AddData(jsonData string, time time.Time, clientUUID string, staffId int, cafeId int) error {
	query := `insert into statistics_table (jsonData, time, clientUUID, staffId, cafeId) VALUES ($1,$2,$3,$4,$5)`
	_, err := p.Conn.Exec(query, jsonData, time, clientUUID, staffId, cafeId)
	return err
}

func NewPostgresStatisticsRepository(conn *sqlx.DB) statistics.Repository {
	return &postgresStatisticsRepository{
		Conn: conn,
	}
}
