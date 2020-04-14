package repository

import (
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/pkg/hasher"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStaffRepository struct {
	conn *sqlx.DB
}

func NewPostgresStaffRepository(conn *sqlx.DB) PostgresStaffRepository {
	cafeStorage := PostgresStaffRepository{conn}
	return cafeStorage
}

func (p *PostgresStaffRepository) Add(ctx context.Context, st models.Staff) (models.Staff, error) {
	query := `INSERT into staff(
                  name, 
                  email, 
                  password,
                  editedat,
                  photo, 
                  isowner, 
                  cafeid, 
                  position) 
                  VALUES ($1,$2,$3,$4,$5,$6,$7,$8) 
                  RETURNING StaffID, Name, Email, EditedAt, Photo, IsOwner, CafeId, Position`

	var dbStaff models.Staff
	hashedPassword, err := hasher.HashAndSalt(nil, st.Password)
	if err != nil {
		return dbStaff, err
	}
	err = p.conn.GetContext(ctx, &dbStaff, query, st.Name, st.Email, hashedPassword, st.EditedAt, st.Photo, st.IsOwner, st.CafeId, st.Position)
	return dbStaff, err
}

func (p *PostgresStaffRepository) GetByEmail(ctx context.Context,
	email string) (models.Staff, error) {

	query := `SELECT StaffID, Name, Email, EditedAt, Photo, IsOwner, CafeId, Position, Password FROM Staff WHERE email=$1`

	var dbStaff models.Staff
	err := p.conn.GetContext(ctx, &dbStaff, query, email)

	return dbStaff, err
}

func (p *PostgresStaffRepository) GetByID(ctx context.Context, id int) (models.Staff, error) {
	query := `SELECT StaffID, Name, Email, EditedAt, Photo, IsOwner, CafeId, Position FROM Staff WHERE StaffID=$1`

	var dbStaff models.Staff
	err := p.conn.GetContext(ctx, &dbStaff, query, id)

	if err != nil {
		return models.Staff{}, err
	}

	return dbStaff, nil
}

func (p *PostgresStaffRepository) Update(ctx context.Context, newStaff models.SafeStaff) (models.SafeStaff, error) {
	query := `UPDATE Staff SET 
                 name=$1,
                 email=$2,
                 editedat=$3,
                 photo=$4,
                 position=$5 WHERE staffid = $6
				RETURNING StaffID, Name, Email, EditedAt, Photo, IsOwner, CafeId, Position`

	var dbStaff models.SafeStaff
	err := p.conn.GetContext(ctx, &dbStaff, query, newStaff.Name, newStaff.Email, newStaff.EditedAt,
		newStaff.Photo, newStaff.Position, newStaff.StaffID)

	return dbStaff, err
}

func (p *PostgresStaffRepository) AddUuid(ctx context.Context, uuid string, id int) error {
	query := `INSERT into UuidCafeRepository(uuid, cafeId) VALUES ($1,$2)`
	_, err := p.conn.ExecContext(ctx, query, uuid, id)
	return err
}

func (p *PostgresStaffRepository) CheckIsOwner(ctx context.Context, staffId int) (bool, error) {
	staff, err := p.GetByID(ctx, staffId)
	if err != nil {
		return false, err
	}
	return staff.IsOwner, nil
}

func (p *PostgresStaffRepository) DeleteUuid(ctx context.Context, uuid string) error {
	query := `DELETE FROM UuidCafeRepository WHERE uuid=$1`
	_, err := p.conn.ExecContext(ctx, query, uuid)
	return err
}

func (p *PostgresStaffRepository) GetCafeId(ctx context.Context, uuid string) (id int, err error) {
	query := `SELECT cafeid FROM uuidcaferepository WHERE uuid=$1`
	err = p.conn.GetContext(ctx, &id, query, uuid)
	return id, err
}

func (p *PostgresStaffRepository) GetStaffListByOwnerId(ctx context.Context, ownerId int) (map[string][]models.StaffByOwnerResponse, error) {
	var data []models.StaffByOwnerResponse
	query := `SELECT 
       cafe.cafeid,
       cafe.cafename,
       s.staffid,
       s.photo,
       s.name,
       s.position from cafe left join staff s on cafe.cafeid = s.cafeid where cafe.staffid=$1 ORDER BY cafe.cafeid
`
	err := p.conn.SelectContext(ctx, &data, query, ownerId)
	if err != nil {
		emptMap := make(map[string][]models.StaffByOwnerResponse)
		return emptMap, err
	}
	return addCafeToList(data), err

}

func (p *PostgresStaffRepository) DeleteStaff(ctx context.Context, staffId int) error {
	query := `DELETE FROM staff where staffid=$1`
	_, err := p.conn.ExecContext(ctx, query, staffId)
	return err
}

func addCafeToList(staffList []models.StaffByOwnerResponse) map[string][]models.StaffByOwnerResponse {
	result := make(map[string][]models.StaffByOwnerResponse)
	for _, staff := range staffList {
		key := staff.CafeId + "," + staff.CafeName
		if staff.StaffId == nil {
			result[key] = nil
			continue
		}
		result[key] = append(result[key], staff)
	}
	return result
}
