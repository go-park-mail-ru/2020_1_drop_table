package repository

import (
	"2020_1_drop_table/internal/app/customer/models"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	type addTestCase struct {
		Customer models.Customer
		err      error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var customer models.Customer
	err = faker.FakeData(&customer)
	assert.NoError(t, err)
	customer.SurveyResult = "{}"

	columnNames := []string{
		"customerid",
		"cafeid",
		"type",
		"points",
		"surveyresult",
	}

	query := `INSERT INTO Customer(
    CafeID, 
    Type,
	Points,
    surveyresult) 
	VALUES ($1,$2,$3,$4) 
	RETURNING *`

	testCases := []addTestCase{
		//Test OK
		{
			Customer: customer,
			err:      nil,
		},
		//Test err
		{
			Customer: models.Customer{},
			err:      sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		row := []driver.Value{testCase.Customer.CustomerID, testCase.Customer.CafeID,
			testCase.Customer.Type, testCase.Customer.Points, "{}"}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(row...)
			// from 1st to delete id
			mock.ExpectQuery(query).WithArgs(row[1:]...).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(row[1:]...).WillReturnError(testCase.err)
		}
		rep := NewPostgresCustomerRepository(sqlxDB)

		passObj, err := rep.Add(context.Background(), testCase.Customer)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.Customer, passObj, message)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestSetLoyaltyPoints(t *testing.T) {
	type addTestCase struct {
		Customer models.Customer
		err      error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var customer models.Customer
	err = faker.FakeData(&customer)
	assert.NoError(t, err)
	customer.SurveyResult = "{}"

	columnNames := []string{
		"customerid",
		"cafeid",
		"type",
		"points",
		"surveyresult",
	}

	query := `UPDATE Customer SET Points=$1 WHERE CustomerID=$2 RETURNING *`

	testCases := []addTestCase{
		//Test OK
		{
			Customer: customer,
			err:      nil,
		},
		//Test err
		{
			Customer: models.Customer{},
			err:      sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		row := []driver.Value{testCase.Customer.CustomerID, testCase.Customer.CafeID,
			testCase.Customer.Type, testCase.Customer.Points, "{}"}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(row...)
			// from 1st to delete id
			mock.ExpectQuery(query).WithArgs(testCase.Customer.Points,
				testCase.Customer.CustomerID).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(testCase.Customer.Points,
				testCase.Customer.CustomerID).WillReturnError(testCase.err)
		}
		rep := NewPostgresCustomerRepository(sqlxDB)

		passObj, err := rep.SetLoyaltyPoints(context.Background(), testCase.Customer.Points,
			testCase.Customer.CustomerID)
		assert.Equal(t, testCase.err, err)
		if err == nil {
			assert.Equal(t, testCase.Customer, passObj, message)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByID(t *testing.T) {
	type addTestCase struct {
		Customer models.Customer
		err      error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var customer models.Customer
	err = faker.FakeData(&customer)
	assert.NoError(t, err)
	customer.SurveyResult = "{}"

	columnNames := []string{
		"customerid",
		"cafeid",
		"type",
		"points",
		"surveyresult",
	}

	query := `SELECT * FROM Customer WHERE CustomerID=$1`

	testCases := []addTestCase{
		//Test OK
		{
			Customer: customer,
			err:      nil,
		},
		//Test err
		{
			Customer: models.Customer{},
			err:      sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		row := []driver.Value{testCase.Customer.CustomerID, testCase.Customer.CafeID,
			testCase.Customer.Type, testCase.Customer.Points, "{}"}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(row...)
			// from 1st to delete id
			mock.ExpectQuery(query).WithArgs(testCase.Customer.CustomerID).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(testCase.Customer.CustomerID).WillReturnError(testCase.err)
		}
		rep := NewPostgresCustomerRepository(sqlxDB)

		passObj, err := rep.GetByID(context.Background(), testCase.Customer.CustomerID)
		assert.Equal(t, testCase.err, err)
		if err == nil {
			assert.Equal(t, testCase.Customer, passObj, message)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestDeleteByID(t *testing.T) {
	type addTestCase struct {
		Customer models.Customer
		err      error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var customer models.Customer
	err = faker.FakeData(&customer)
	assert.NoError(t, err)
	customer.SurveyResult = "{}"

	query := `DELETE FROM Customer WHERE CustomerID=$1`

	testCases := []addTestCase{
		//Test OK
		{
			Customer: customer,
			err:      nil,
		},
		//Test err
		{
			Customer: models.Customer{},
			err:      sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		if testCase.err == nil {
			mock.ExpectExec(query).WithArgs(
				testCase.Customer.CustomerID).WillReturnResult(sqlmock.NewResult(0, 0))
		} else {
			mock.ExpectExec(query).WithArgs(testCase.Customer.CustomerID).WillReturnError(testCase.err)
		}
		rep := NewPostgresCustomerRepository(sqlxDB)

		err := rep.DeleteByID(context.Background(), testCase.Customer.CustomerID)
		assert.Equal(t, testCase.err, err, message)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
