package repository

import (
	"2020_1_drop_table/internal/app/apple_passkit"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"2020_1_drop_table/internal/pkg/apple_pass_generator/meta"
	"context"
	"database/sql"
	"encoding/json"
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
    CafeID,
    Type,        
	LoyaltyInfo, 
	published,   
	Design, 
	Icon, 
	Icon2x, 
	Logo, 
	Logo2x, 
	Strip, 
	Strip2x) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) 
	RETURNING *`

	var dbApplePass models.ApplePassDB
	err := p.Conn.GetContext(ctx, &dbApplePass, query, ap.CafeID, ap.Type, ap.LoyaltyInfo, ap.Published, ap.Design,
		ap.Icon, ap.Icon2x, ap.Logo, ap.Logo2x, ap.Strip, ap.Strip2x)

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

func (p *postgresApplePassRepository) GetPassByCafeID(ctx context.Context, cafeID int, Type string, published bool) (models.ApplePassDB, error) {
	query := `SELECT * FROM ApplePass WHERE CafeID=$1 AND Type=$2 AND published=$3`

	var dbApplePass models.ApplePassDB
	err := p.Conn.GetContext(ctx, &dbApplePass, query, cafeID, Type, published)

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
	Strip2x=NotEmpty($7, Strip2x),
    LoyaltyInfo=NotEmpty($8, LoyaltyInfo)
	WHERE CafeID=$9 AND Type=$10 AND published=$11`

	_, err := p.Conn.ExecContext(ctx, query, newApplePass.Design, newApplePass.Icon,
		newApplePass.Icon2x, newApplePass.Logo, newApplePass.Logo2x,
		newApplePass.Strip, newApplePass.Strip2x, newApplePass.LoyaltyInfo, newApplePass.CafeID,
		newApplePass.Type, newApplePass.Published)

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

func (p *postgresApplePassRepository) createMeta(ctx context.Context, cafeID int) error {
	query := `INSERT INTO ApplePassMeta(
	CafeID,
    meta)
	VALUES ($1, $2)`
	emptyMeta, err := json.Marshal(meta.EmptyMeta)
	if err != nil {
		return err
	}

	_, err = p.Conn.ExecContext(ctx, query, cafeID, emptyMeta)
	if err != nil {
		return err
	}

	return err
}

func (p *postgresApplePassRepository) UpdateMeta(ctx context.Context, cafeID int, meta []byte) error {
	query := `UPDATE ApplePassMeta 
    SET meta=$1
	WHERE CafeID=$2`

	_, err := p.Conn.ExecContext(ctx, query, meta, cafeID)

	if err != nil {
		if err == sql.ErrNoRows {
			return p.createMeta(ctx, cafeID)
		}
		return err
	}

	return err
}

func (p *postgresApplePassRepository) GetMeta(ctx context.Context,
	cafeID int) (applePassMeta models.ApplePassMeta, err error) {

	query := `SELECT meta FROM ApplePassMeta WHERE CafeID=$1`

	var metaJson []byte
	err = p.Conn.GetContext(ctx, &metaJson, query, cafeID)
	if err != nil {
		if err == sql.ErrNoRows {
			applePassMeta.CafeID = cafeID
			applePassMeta.Meta = meta.EmptyMeta

			return applePassMeta, p.createMeta(ctx, cafeID)
		}
		return models.ApplePassMeta{}, err
	}

	applePassMeta.CafeID = cafeID
	//TODO can be bug with EasyJSOn not tested
	err = applePassMeta.UnmarshalJSON(metaJson)
	if err != nil {
		return models.ApplePassMeta{}, err
	}

	return applePassMeta, err
}
