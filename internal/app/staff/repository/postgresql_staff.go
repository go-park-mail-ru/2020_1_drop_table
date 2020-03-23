package repository

import (
	"2020_1_drop_table/internal/app/staff/models"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgresStaffRepository struct {
	Conn *sqlx.DB
}

func NewPostgresStaffRepository(conn *sqlx.DB) postgresStaffRepository {
	cafeStorage := postgresStaffRepository{conn}
	return cafeStorage
}

func (p *postgresStaffRepository) Add(ctx context.Context, st models.Staff) (models.Staff, error) {
	query := `INSERT into staff(name, email, password, editedat, photo, isowner) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *`
	var dbStaff models.Staff
	err := p.Conn.GetContext(ctx, &dbStaff, query, st.Name, st.Email, st.Password, st.EditedAt, st.Photo, st.IsOwner)
	return dbStaff, err
}

func (p *postgresStaffRepository) GetByEmailAndPassword(ctx context.Context,
	email string, password string) (models.Staff, error) {

	query := `SELECT * FROM Staff WHERE password=$1 AND email=$2`

	var dbStaff models.Staff
	err := p.Conn.GetContext(ctx, &dbStaff, query, password, email)

	return dbStaff, err
}

func (p *postgresStaffRepository) GetByID(ctx context.Context, id int) (models.Staff, error) {
	query := `SELECT * FROM Staff WHERE StaffID=$1`

	var dbStaff models.Staff
	err := p.Conn.GetContext(ctx, &dbStaff, query, id)

	if err != nil {
		return models.Staff{}, err
	}

	return dbStaff, nil
}

func (p *postgresStaffRepository) Update(ctx context.Context, newStaff models.SafeStaff) error {
	query := `UPDATE Staff SET name=$1,email=$2,editedat=$3,photo=$4 WHERE staffid = $5`
	_, err := p.Conn.ExecContext(ctx, query, newStaff.Name, newStaff.Email, newStaff.EditedAt,
		newStaff.Photo, newStaff.StaffID)

	return err
}

func (p *postgresStaffRepository) AddUuid(ctx context.Context, uuid string, id int) error {
	query := `INSERT into UuidCafeRepository(uuid, cafeId) VALUES ($1,$2)`
	_, err := p.Conn.ExecContext(ctx, query, uuid, id)
	return err
}

func (p *postgresStaffRepository) CheckIsOwner(ctx context.Context, staffId int) (bool, error) {
	staff, err := p.GetByID(ctx, staffId)
	if err != nil {
		return false, err
	}
	return staff.IsOwner == true, nil
}

func (p *postgresStaffRepository) DeleteUuid(ctx context.Context, uuid string) error {
	query := `DELETE FROM UuidCafeRepository WHERE uuid=$1`
	_, err := p.Conn.ExecContext(ctx, query, uuid)
	return err

}
