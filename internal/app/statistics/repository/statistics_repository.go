package repository

import (
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	"2020_1_drop_table/internal/app/statistics"
	"2020_1_drop_table/internal/app/statistics/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

type postgresStatisticsRepository struct {
	Conn *sqlx.DB
}

func generateWhereStatement(cafeList []cafeModels.Cafe) string {
	query := `where(`
	for _, cafe := range cafeList {
		query = query + " cafeID=" + strconv.Itoa(cafe.CafeID) + " or "
	}
	query = query[:len(query)-3] + ")"
	return query
}

func generateDateTruncStatement(typ string) string {
	if typ == "MONTH" {
		return `date_trunc('MONTH', time)`
	}
	return `date_trunc('DAY', time)`
}

func generateBetweenStatement(dateTruncStatement string, since string, to string) string {
	return fmt.Sprintf(`(%s between '%s' and '%s')`, dateTruncStatement, since, to)
}

func (p postgresStatisticsRepository) GetGraphsDataFromRepo(ctx context.Context, cafeList []cafeModels.Cafe, typ string, since string, to string) ([]models.StatisticsGraphRawStruct, error) {
	whereStatement := generateWhereStatement(cafeList)
	dateTrunc := generateDateTruncStatement(typ)
	betweenStatement := generateBetweenStatement(dateTrunc, since, to)
	query := fmt.Sprintf(`SELECT count(*), %s as date, cafeID, staffID
from statistics_table
%s and %s
group by %s, cafeID, staffID
order by cafeID, %s`, dateTrunc, whereStatement, betweenStatement, dateTrunc, dateTrunc)
	fmt.Println(query)
	var res []models.StatisticsGraphRawStruct
	err := p.Conn.SelectContext(ctx, &res, query)
	return res, err

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
