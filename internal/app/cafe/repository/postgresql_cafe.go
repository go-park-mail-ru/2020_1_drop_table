package repository

import (
	"2020_1_drop_table/internal/app/cafe/models"
	"context"
	"github.com/jmoiron/sqlx"
)

type postgresCafeRepository struct {
	Conn *sqlx.DB
}

func NewPostgresCafeRepository(conn *sqlx.DB) postgresCafeRepository {
	cafeStorage := postgresCafeRepository{conn}
	return cafeStorage
}

func (p *postgresCafeRepository) Add(ctx context.Context, ca models.Cafe) (models.Cafe, error) {
	query := `INSERT INTO Cafe(
	Name, 
	Address, 
	Description, 
	StaffID, 
	OpenTime, 
	CloseTime, 
	Photo) 
	VALUES ($1,$2,$3,$4,$5,$6,$7) 
	RETURNING *`

	var dbCafe models.Cafe
	err := p.Conn.GetContext(ctx, &dbCafe, query, ca.Name, ca.Address,
		ca.Description, ca.StaffID, ca.OpenTime, ca.CloseTime, ca.Photo)

	return dbCafe, err
}

func (p *postgresCafeRepository) GetByID(ctx context.Context, id int) (models.Cafe, error) {
	query := `SELECT * FROM Cafe WHERE CafeID=$1`

	var dbStaff models.Cafe
	err := p.Conn.GetContext(ctx, &dbStaff, query, id)

	if err != nil {
		return models.Cafe{}, err
	}

	return dbStaff, nil
}

func (p *postgresCafeRepository) GetByOwnerID(ctx context.Context, staffID int) ([]models.Cafe, error) {
	query := `SELECT * FROM Cafe WHERE StaffID=$1 ORDER BY CafeID`

	var cafes []models.Cafe
	err := p.Conn.SelectContext(ctx, &cafes, query, staffID)

	if err != nil {
		return make([]models.Cafe, 0), err
	}

	return cafes, nil
}

func (p *postgresCafeRepository) Update(ctx context.Context, newCafe models.Cafe) error {
	query := `UPDATE Cafe SET 
	Name=$1, 
	Address=$2, 
	Description=$3, 
	StaffID=$4, 
	OpenTime=$5, 
	CloseTime=$6, 
	Photo=$7 
	WHERE CafeID=$8`

	_, err := p.Conn.ExecContext(ctx, query, newCafe.Name, newCafe.Address, newCafe.Description,
		newCafe.StaffID, newCafe.OpenTime, newCafe.CloseTime, newCafe.Photo, newCafe.CafeID)

	return err
}
