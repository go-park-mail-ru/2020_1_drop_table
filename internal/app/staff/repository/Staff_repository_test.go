package repository

import (
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/pkg/hasher"
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func addGetByEmailSupport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"staffid", "name", "email", "password", "editedat", "photo", "isowner", "cafeid", "position"}).
		AddRow(1, "test", "valid@valid.ru", "123", time.Now().UTC(), "photo", true, 0, "position")
	mock.ExpectQuery("SELECT * FROM Staff WHERE email=$1").WithArgs("valid@valid.ru").WillReturnRows(rows)
	mock.ExpectQuery("SELECT * FROM Staff WHERE email=$1").WithArgs("notexist").WillReturnError(sql.ErrNoRows)
}

func addUpdateSupport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"staffid", "name", "email", "editedat", "photo", "isowner", "cafeid", "position"}).
		AddRow(1, "test", "valid@valid.ru", time.Now().UTC(), "photo", true, 0, "position")
	mock.ExpectQuery(`UPDATE Staff SET name=$1,email=$2,editedat=$3,photo=$4 WHERE staffid = $5 RETURNING StaffID, Name, Email,
			  EditedAt, Photo, IsOwner, CafeId, Position`).WithArgs("test", "valid@valid.ru", time.Time{}, "", 1).WillReturnRows(rows)

}

func getDataBase() (*sqlx.DB, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	addGetByEmailSupport(mock)
	addGetByIdSupport(mock)
	addAddSupport(mock)
	addUpdateSupport(mock)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, err
}

func addGetByIdSupport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"staffid", "name", "email", "password", "editedat", "photo", "isowner", "cafeid", "position"}).
		AddRow(2, "tesxct", "valid@valid.ru", "123", time.Now().UTC(), "photo", true, 0, "position")
	mock.ExpectQuery("SELECT * FROM Staff WHERE StaffID=$1").WithArgs(2).WillReturnRows(rows)
	mock.ExpectQuery("SELECT * FROM Staff WHERE StaffID=$1").WithArgs(-228).WillReturnError(sql.ErrNoRows)
}

func addAddSupport(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"staffid", "name", "email", "password", "editedat", "photo", "isowner", "cafeid", "position"}).
		AddRow(2, "test", "valid@valid.ru", "123", time.Now().UTC(), "photo", true, 0, "position")
	pass, _ := hasher.HashAndSalt(nil, "123")
	mock.ExpectQuery("INSERT into staff(name, email, password, editedat, photo, isowner, cafeid, position) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *").WithArgs("test", "valid@valid.ru", pass, time.Time{}, "", false, 0, "").WillReturnRows(rows)
}

func getEmptyDb() (*sqlx.DB, sqlmock.Sqlmock) {

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	con := sqlx.NewDb(db, "sqlmock")
	return con, mock
}

func TestAdd(t *testing.T) {
	st := models.Staff{
		StaffID:  2,
		Name:     "test",
		Email:    "valid@valid.ru",
		Password: "123",
	}
	con, err := getDataBase()
	rep := NewPostgresStaffRepository(con)
	_, err = rep.Add(context.TODO(), st)
	assert.NotNil(t, err)

}

func TestGetByEmail(t *testing.T) {
	con, err := getDataBase()
	resUser := models.Staff{
		StaffID:  1,
		Name:     "test",
		Email:    "valid@valid.ru",
		Password: "123",
		EditedAt: time.Now().UTC(),
		Photo:    "photo",
		IsOwner:  true,
		CafeId:   0,
		Position: "position",
	}
	rep := NewPostgresStaffRepository(con)
	res, err := rep.GetByEmail(context.TODO(), "valid@valid.ru")
	assert.Nil(t, err)
	assert.Equal(t, resUser.Email, res.Email)
	assert.Equal(t, resUser.Password, res.Password)
	res, err = rep.GetByEmail(context.TODO(), "notexist")
	assert.NotNil(t, err)
}

func TestGetById(t *testing.T) {
	con, err := getDataBase()
	rep := NewPostgresStaffRepository(con)
	_, err = rep.GetByID(context.TODO(), -228)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	con, mock := getEmptyDb()
	addUpdateSupport(mock)
	rep := NewPostgresStaffRepository(con)
	resUser := models.SafeStaff{
		StaffID:  1,
		Name:     "test",
		Email:    "valid@valid.ru",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  true,
		CafeId:   0,
		Position: "position",
	}
	res, err := rep.Update(context.TODO(), resUser)
	assert.Nil(t, err)
	assert.Equal(t, res.Email, resUser.Email)
}

func TestAddUuid(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	con := sqlx.NewDb(db, "sqlmock")
	mock.ExpectQuery("INSERT into UuidCafeRepository(uuid, cafeId) VALUES ($1,$2)").WithArgs("asdasdasdasd", -1).WillReturnError(nil)
	rep := NewPostgresStaffRepository(con)
	err := rep.AddUuid(context.TODO(), "asdasdasdasd", -1)
	assert.NotNil(t, err)
}
