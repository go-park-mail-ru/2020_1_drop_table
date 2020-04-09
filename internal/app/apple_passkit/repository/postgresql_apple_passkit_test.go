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
	"strconv"
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

func TestUpdate(t *testing.T) {
	type addTestCase struct {
		inputPass passKitModels.ApplePassDB
		err       error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var inputPass passKitModels.ApplePassDB
	err = faker.FakeData(&inputPass)
	assert.NoError(t, err)

	query := `UPDATE ApplePass SET  
	 Design=NotEmpty($1, Design),
	 Icon=NotEmpty($2, Icon),
	 Icon2x=NotEmpty($3, Icon2x),
	 Logo=NotEmpty($4, Logo),
	 Logo2x=NotEmpty($5, Logo2x),
	 Strip=NotEmpty($6, Strip),
	 Strip2x=NotEmpty($7, Strip2x)
	 WHERE ApplePassID=$8`

	testCases := []addTestCase{
		//Test OK
		{
			inputPass: inputPass,
			err:       nil,
		},
		//Test not found
		{
			inputPass: passKitModels.ApplePassDB{},
			err:       sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		args := []driver.Value{testCase.inputPass.ApplePassID, testCase.inputPass.Design,
			testCase.inputPass.Icon, testCase.inputPass.Icon2x, testCase.inputPass.Logo,
			testCase.inputPass.Logo2x, testCase.inputPass.Strip, testCase.inputPass.Strip2x}
		req := mock.ExpectExec(query).WithArgs(append(args[1:], args[0])...)
		if testCase.err == nil {
			// append is needed to make cafeID last param
			// till the second before end to delete apple passes IDs
			req.WillReturnResult(sqlmock.NewResult(0, 0))
		} else {
			req.WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		err := rep.Update(context.Background(), testCase.inputPass)
		assert.Equal(t, testCase.err, err, message)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestUpdateDesign(t *testing.T) {
	type addTestCase struct {
		id     int
		design string
		err    error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var design string
	err = faker.FakeData(&design)
	assert.NoError(t, err)

	var id int
	err = faker.FakeData(&id)
	assert.NoError(t, err)

	query := `UPDATE ApplePass SET Design=$1 WHERE ApplePassID=$2`

	testCases := []addTestCase{
		//Test OK
		{
			id:     id,
			design: design,
			err:    nil,
		},
		//Test not found
		{
			id:     -1,
			design: design,
			err:    sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		args := []driver.Value{testCase.design, testCase.id}

		req := mock.ExpectExec(query).WithArgs(args...)

		if testCase.err == nil {
			// append is needed to make cafeID last param
			// till the second before end to delete apple passes IDs
			req.WillReturnResult(sqlmock.NewResult(0, 0))
		} else {
			req.WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		err := rep.UpdateDesign(context.Background(), testCase.design, testCase.id)
		assert.Equal(t, testCase.err, err, message)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestDelete(t *testing.T) {
	type addTestCase struct {
		id  int
		err error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var design string
	err = faker.FakeData(&design)
	assert.NoError(t, err)

	var id int
	err = faker.FakeData(&id)
	assert.NoError(t, err)

	query := `DELETE FROM ApplePass WHERE ApplePassID=$2`

	testCases := []addTestCase{
		//Test OK
		{
			id:  id,
			err: nil,
		},
		//Test not found
		{
			id:  -1,
			err: sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		args := []driver.Value{testCase.id}

		req := mock.ExpectExec(query).WithArgs(args...)

		if testCase.err == nil {
			// append is needed to make cafeID last param
			// till the second before end to delete apple passes IDs
			req.WillReturnResult(sqlmock.NewResult(0, 0))
		} else {
			req.WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		err := rep.Delete(context.Background(), testCase.id)
		assert.Equal(t, testCase.err, err, message)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestUpdateMeta(t *testing.T) {
	type addTestCase struct {
		cafeID     int
		outputMeta passKitModels.ApplePassMeta
		query      string
		err        error
		finalErr   error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var cafeID int
	err = faker.FakeData(&cafeID)
	assert.NoError(t, err)

	var outputMeta passKitModels.ApplePassMeta
	err = faker.FakeData(&outputMeta)
	assert.NoError(t, err)
	outputMeta.CafeID = strconv.Itoa(cafeID)

	columnNames := []string{
		"applepassmetaid",
		"cafeid",
		"passescount",
	}

	query := `UPDATE ApplePassMeta 
    SET PassesCount = PassesCount + 1
	WHERE CafeID=$1 
	RETURNING *`

	query2 := `INSERT INTO ApplePassMeta(
	CafeID,
	PassesCount)
	VALUES ($1,1)
	RETURNING *`

	testCases := []addTestCase{
		//Test OK
		{
			cafeID:     cafeID,
			outputMeta: outputMeta,
			query:      "",
			err:        nil,
		},
		//Test error
		{
			cafeID:     cafeID + 1,
			outputMeta: outputMeta,
			query:      query2,
			err:        sql.ErrNoRows,
			finalErr:   nil,
		},
		//Test no cafe
		{
			cafeID:     cafeID + 1,
			outputMeta: outputMeta,
			query:      query2,
			err:        sql.ErrNoRows,
			finalErr:   sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		rows := []driver.Value{testCase.outputMeta.ApplePassMetaID, testCase.outputMeta.CafeID,
			testCase.outputMeta.PassesCount}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(rows...)
			// from 1st to delete id
			mock.ExpectQuery(query).WithArgs(testCase.cafeID).WillReturnRows(rows)
		} else {
			rows := sqlmock.NewRows(columnNames).AddRow(rows...)
			mock.ExpectQuery(query).WithArgs(testCase.cafeID).WillReturnError(testCase.err)

			if testCase.finalErr == nil {
				mock.ExpectQuery(testCase.query).WithArgs(testCase.cafeID).WillReturnRows(rows)
			} else {
				mock.ExpectQuery(testCase.query).WithArgs(testCase.cafeID).WillReturnError(testCase.finalErr)
			}
		}
		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		metaObj, err := rep.UpdateMeta(context.Background(), testCase.cafeID)
		assert.Equal(t, testCase.finalErr, err, message)
		if err == nil {
			assert.Equal(t, testCase.outputMeta, metaObj, message)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
