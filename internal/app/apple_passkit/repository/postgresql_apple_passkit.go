package repository

import (
	"2020_1_drop_table/internal/app/apple_passkit"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"context"
	"github.com/jmoiron/sqlx"
)

type postgresApplePassRepository struct {
	Conn *sqlx.DB
}

func NewPostgresApplePassRepository(conn *sqlx.DB) apple_passkit.Repository {
	return &postgresApplePassRepository{
		Conn: conn,
	}
}

func (p *postgresApplePassRepository) Add(ctx context.Context, ap models.ApplePassDB) (models.ApplePassDB, error) {
	query := `INSERT INTO ApplePass(
	Design, 
	Icon, 
	Icon2x, 
	Logo, 
	Logo2x, 
	Strip, 
	Strip2x) 
	VALUES ($1,$2,$3,$4,$5,$6,$7) 
	RETURNING *`

	var dbApplePass models.ApplePassDB
	err := p.Conn.GetContext(ctx, &dbApplePass, query, ap.Design, ap.Icon, ap.Icon2x, ap.Logo,
		ap.Logo2x, ap.Strip, ap.Strip2x)

	if err != nil {
		return models.ApplePassDB{}, err
	}

	return dbApplePass, err
}

func (p *postgresApplePassRepository) GetPassByID(ctx context.Context, id int) (models.ApplePassDB, error) {
	query := `SELECT * FROM ApplePass WHERE ApplePassID=$1`

	var dbApplePass models.ApplePassDB
	err := p.Conn.GetContext(ctx, &dbApplePass, query, id)

	if err != nil {
		return models.ApplePassDB{}, err
	}

	return dbApplePass, err
}

func (p *postgresApplePassRepository) GetDesignByID(ctx context.Context, id int) (models.ApplePassDB, error) {
	query := `SELECT Design FROM ApplePass WHERE ApplePassID=$1`

	var dbApplePass models.ApplePassDB
	err := p.Conn.GetContext(ctx, &dbApplePass, query, id)

	if err != nil {
		return models.ApplePassDB{}, err
	}

	return dbApplePass, err
}

func (p *postgresApplePassRepository) Update(ctx context.Context, newApplePass models.ApplePassDB) error {
	query := `UPDATE ApplePass SET 
	Design=$1, 
	Icon=$2, 
	Icon2x=$3, 
	Logo=$4, 
	Logo2x=$5, 
	Strip=$6, 
	Strip2x=$7 
	WHERE ApplePassID=$8`

	_, err := p.Conn.ExecContext(ctx, query, newApplePass.Design, newApplePass.Icon, newApplePass.Icon2x,
		newApplePass.Logo, newApplePass.Logo2x, newApplePass.Strip, newApplePass.Strip2x,
		newApplePass.ApplePassID)

	return err
}

func (p *postgresApplePassRepository) UpdateDesign(ctx context.Context, Design string, id int) error {
	query := `UPDATE ApplePass SET Design=$1 WHERE ApplePassID=$2`

	_, err := p.Conn.ExecContext(ctx, query, Design, id)

	return err
}

func (p *postgresApplePassRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM ApplePass WHERE ApplePassID=$2`

	_, err := p.Conn.ExecContext(ctx, query, id)

	return err
}
