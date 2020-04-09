package repository_test

import (
	passKitModels "2020_1_drop_table/internal/app/apple_passkit/models"
	"2020_1_drop_table/internal/app/apple_passkit/repository"
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
		inputPass  passKitModels.ApplePassDB
		outputPass passKitModels.ApplePassDB
		err        error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var outputPass passKitModels.ApplePassDB
	err = faker.FakeData(&outputPass)
	assert.NoError(t, err)

	inputPass := outputPass
	inputPass.ApplePassID = 0

	columnNames := []string{
		"applepassid",
		"design",
		"icon",
		"icon2x",
		"logo",
		"logo2x",
		"strip",
		"strip2x",
	}

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

	testCases := []addTestCase{
		//Test OK
		{
			inputPass:  inputPass,
			outputPass: outputPass,
			err:        nil,
		},
		//Test error
		{
			inputPass:  passKitModels.ApplePassDB{},
			outputPass: passKitModels.ApplePassDB{},
			err:        sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		args := []driver.Value{testCase.outputPass.ApplePassID, testCase.outputPass.Design,
			testCase.outputPass.Icon, testCase.outputPass.Icon2x, testCase.outputPass.Logo,
			testCase.outputPass.Logo2x, testCase.outputPass.Strip, testCase.outputPass.Strip2x}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(args...)
			// from 1st to delete id
			mock.ExpectQuery(query).WithArgs(args[1:]...).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(args[1:]...).WillReturnError(testCase.err)
		}
		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		passObj, err := rep.Add(context.Background(), testCase.inputPass)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.outputPass, passObj, message)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByID(t *testing.T) {
	type addTestCase struct {
		applePassID int
		outputPass  passKitModels.ApplePassDB
		err         error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var outputPass passKitModels.ApplePassDB
	err = faker.FakeData(&outputPass)
	assert.NoError(t, err)

	inputPass := outputPass
	inputPass.ApplePassID = 0

	columnNames := []string{
		"applepassid",
		"design",
		"icon",
		"icon2x",
		"logo",
		"logo2x",
		"strip",
		"strip2x",
	}

	query := `SELECT * FROM ApplePass WHERE ApplePassID=$1`

	testCases := []addTestCase{
		//Test OK
		{
			applePassID: outputPass.ApplePassID,
			outputPass:  outputPass,
			err:         nil,
		},
		//Test not found
		{
			applePassID: -1,
			outputPass:  passKitModels.ApplePassDB{},
			err:         sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		row := []driver.Value{testCase.outputPass.ApplePassID, testCase.outputPass.Design,
			testCase.outputPass.Icon, testCase.outputPass.Icon2x, testCase.outputPass.Logo,
			testCase.outputPass.Logo2x, testCase.outputPass.Strip, testCase.outputPass.Strip2x}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(row...)
			// from 1st to delete id
			mock.ExpectQuery(query).WithArgs(testCase.applePassID).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(testCase.applePassID).WillReturnError(testCase.err)
		}
		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		passObj, err := rep.GetPassByID(context.Background(), testCase.applePassID)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.outputPass, passObj, message)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
