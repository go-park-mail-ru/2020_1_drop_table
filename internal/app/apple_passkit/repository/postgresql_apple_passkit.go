package repository

import (
	"2020_1_drop_table/internal/app/apple_passkit"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"context"
	"database/sql"
	"fmt"
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
		fmt.Println("HERE")
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

func (p *postgresApplePassRepository) Update(ctx context.Context, newApplePass models.ApplePassDB) error {
	query := `UPDATE ApplePass SET  
	Design=NotEmpty($1, Design),
	Icon=NotEmpty($2, Icon),
	Icon2x=NotEmpty($3, Icon2x),
	Logo=NotEmpty($4, Logo),
	Logo2x=NotEmpty($5, Logo2x),
	Strip=NotEmpty($6, Strip),
	Strip2x=NotEmpty($7, Strip2x)
	WHERE ApplePassID=$8`

	_, err := p.Conn.ExecContext(ctx, query, newApplePass.Design, newApplePass.Icon,
		newApplePass.Icon2x, newApplePass.Logo, newApplePass.Logo2x,
		newApplePass.Strip, newApplePass.Strip2x, newApplePass.ApplePassID)

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

func (p *postgresApplePassRepository) createMeta(ctx context.Context, cafeID int) (models.ApplePassMeta, error) {
	query := `INSERT INTO ApplePassMeta(
	CafeID,
	PassesCount)
	VALUES ($1,1)
	RETURNING *`

	var dbApplePass models.ApplePassMeta
	err := p.Conn.GetContext(ctx, &dbApplePass, query, cafeID)

	if err != nil {
		return models.ApplePassMeta{}, err
	}

	return dbApplePass, err
}

func (p *postgresApplePassRepository) UpdateMeta(ctx context.Context, cafeID int) (models.ApplePassMeta, error) {
	query := `UPDATE ApplePassMeta 
    SET PassesCount = PassesCount + 1
	WHERE CafeID=$1 
	RETURNING *`

	var applePassMeta models.ApplePassMeta
	err := p.Conn.GetContext(ctx, &applePassMeta, query, cafeID)

	if err != nil {
		if err == sql.ErrNoRows {
			return p.createMeta(ctx, cafeID)
		}
		return models.ApplePassMeta{}, err
	}

	return applePassMeta, err
}
