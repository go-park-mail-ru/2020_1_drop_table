package test

import (
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/app/staff/repository"
	"context"
	"database/sql"
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

func getDataBase() (*sqlx.DB, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	addGetByEmailSupport(mock)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, err
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
	rep := repository.NewPostgresStaffRepository(con)
	res, err := rep.GetByEmail(context.TODO(), "valid@valid.ru")
	assert.Nil(t, err)
	assert.Equal(t, resUser.Email, res.Email)
	assert.Equal(t, resUser.Password, res.Password)
	res, err = rep.GetByEmail(context.TODO(), "notexist")
	assert.NotNil(t, err)

}
