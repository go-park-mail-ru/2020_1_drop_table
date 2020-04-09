package repository

import (
	"2020_1_drop_table/internal/app/customer/models"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

var columnNames = []string{
	"customerid",
	"cafeid",
	"points",
}

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	rep := NewPostgresCustomerRepository(sqlxDB)
	query := `INSERT INTO Customer(
    CafeID, 
	Points) 
	VALUES ($1,$2) 
	RETURNING *`
	ret := models.Customer{
		CustomerID: "1",
		CafeID:     1,
		Points:     228,
	}
	rows := sqlmock.NewRows(columnNames).AddRow("1", 1, 228)
	mock.ExpectQuery(query).WithArgs(1, 228).WillReturnRows(rows)
	cust, err := rep.Add(context.TODO(), ret)
	assert.Nil(t, err)
	assert.Equal(t, cust, ret)

}

func TestSetLoyalityPoints(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	rep := NewPostgresCustomerRepository(sqlxDB)
	query := `UPDATE Customer SET Points=$1 WHERE CustomerID=$2 RETURNING *`
	ret := models.Customer{
		CustomerID: "1",
		CafeID:     1,
		Points:     228,
	}
	rows := sqlmock.NewRows(columnNames).AddRow("1", 1, 228)
	mock.ExpectQuery(query).WithArgs(228, "1").WillReturnRows(rows)
	cust, err := rep.SetLoyaltyPoints(context.TODO(), ret.Points, ret.CustomerID)
	assert.Nil(t, err)
	assert.Equal(t, cust, ret)
}

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	rep := NewPostgresCustomerRepository(sqlxDB)
	query := `SELECT * FROM Customer WHERE CustomerID=$1`
	ret := models.Customer{
		CustomerID: "1",
		CafeID:     1,
		Points:     228,
	}
	rows := sqlmock.NewRows(columnNames).AddRow("1", 1, 228)
	mock.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)
	cust, err := rep.GetByID(context.TODO(), ret.CustomerID)
	assert.Nil(t, err)
	assert.Equal(t, cust, ret)
}
