package repository

import (
	"2020_1_drop_table/internal/app/cafe"
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type postgresCafeRepository struct {
	Conn *sqlx.DB
}

func (p *postgresCafeRepository) SearchCafes(ctx context.Context, searchBy string, limit int, since int) ([]models.Cafe, error) {
	query := `  SELECT CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str
				FROM cafe
				WHERE CafeName % $1
				   or Address % $1
				   or CafeName LIKE '%' || $1 || '%'
				   or Address  LIKE '%' || $1 || '%'
				limit $2
				offset $3`
	var cafes []models.Cafe
	err := p.Conn.SelectContext(ctx, &cafes, query, searchBy, limit, since)
	return cafes, err

}

func NewPostgresCafeRepository(conn *sqlx.DB) cafe.Repository {
	return &postgresCafeRepository{
		Conn: conn,
	}
}

func GeneratePointToGeo(latitude string, longitude string) string {
	return fmt.Sprintf("SRID=4326;POINT(%s %s)", latitude, longitude)
}
func GeneratePointToGeoWithPoint(point string) string {
	return fmt.Sprintf("SRID=4326;POINT(%s)", point)
}

func (p *postgresCafeRepository) GetCafeSortedByRadius(ctx context.Context, latitude string, longitude string, radius string) ([]models.Cafe, error) {
	point := GeneratePointToGeo(latitude, longitude)
	var resArr []models.Cafe
	query := `SELECT CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str
              FROM cafe where ST_Distance(location::geography, $1::geography)<$2 
              ORDER BY location <-> $1`

	err := p.Conn.SelectContext(ctx, &resArr, query, point, radius)
	return resArr, err
}

func (p *postgresCafeRepository) Add(ctx context.Context, ca models.Cafe) (models.Cafe, error) {
	query := `INSERT INTO Cafe(
	CafeName, 
	Address, 
	Description, 
	StaffID, 
	OpenTime,
	CloseTime, 
	Photo,
    location,
    location_str) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,ST_GeomFromEWKT($8),$9) 
	RETURNING CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str`

	var dbCafe models.Cafe
	postGisPoint := GeneratePointToGeoWithPoint(ca.Location)
	err := p.Conn.GetContext(ctx, &dbCafe, query, ca.CafeName, ca.Address,
		ca.Description, ca.StaffID, ca.OpenTime, ca.CloseTime, ca.Photo, postGisPoint, ca.Location)

	return dbCafe, err
}

func (p *postgresCafeRepository) GetByID(ctx context.Context, id int) (models.Cafe, error) {
	query := `SELECT CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str FROM Cafe WHERE CafeID=$1`

	var dbCafe models.Cafe
	err := p.Conn.GetContext(ctx, &dbCafe, query, id)

	if err != nil {
		return models.Cafe{}, err
	}
	return dbCafe, nil
}

func (p *postgresCafeRepository) GetByOwnerID(ctx context.Context, staffID int) ([]models.Cafe, error) {
	query := `SELECT CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str FROM Cafe WHERE StaffID=$1 ORDER BY CafeID`

	var cafes []models.Cafe
	err := p.Conn.SelectContext(ctx, &cafes, query, staffID)

	if err != nil {
		return make([]models.Cafe, 0), err
	}

	return cafes, nil
}

func (p *postgresCafeRepository) Update(ctx context.Context, newCafe models.Cafe) (models.Cafe, error) {
	query := `UPDATE Cafe SET 
	CafeName=$1, 
	Address=$2, 
	Description=$3, 
	OpenTime=$4, 
	CloseTime=$5, 
	Photo=NotEmpty($6,Photo) 
	WHERE CafeID=$7
	RETURNING CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str`

	var CafeDB models.Cafe

	err := p.Conn.GetContext(ctx, &CafeDB, query, newCafe.CafeName, newCafe.Address, newCafe.Description,
		newCafe.OpenTime, newCafe.CloseTime, newCafe.Photo, newCafe.CafeID)

	return CafeDB, err
}

func (p *postgresCafeRepository) GetAllCafes(ctx context.Context, since int, limit int) ([]models.Cafe, error) {
	query := `SELECT CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str from cafe OFFSET $1 LIMIT $2`
	var CafesList []models.Cafe
	err := p.Conn.SelectContext(ctx, &CafesList, query, since, limit)
	return CafesList, err
}
