package repository_test

import (
	passKitModels "2020_1_drop_table/internal/app/apple_passkit/models"
	"2020_1_drop_table/internal/app/apple_passkit/repository"
	"2020_1_drop_table/internal/pkg/apple_pass_generator/meta"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
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

	assert.NoError(t, err)

	inputPass := outputPass
	inputPass.ApplePassID = 0

	columnNames := []string{
		"applepassid",
		"cafeid",
		"type",
		"loyaltyinfo",
		"published",
		"design",
		"icon",
		"icon2x",
		"logo",
		"logo2x",
		"strip",
		"strip2x",
	}

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

		args := []driver.Value{testCase.outputPass.ApplePassID, testCase.outputPass.CafeID,
			testCase.outputPass.Type, testCase.outputPass.LoyaltyInfo, testCase.outputPass.Published,
			testCase.outputPass.Design, testCase.outputPass.Icon, testCase.outputPass.Icon2x,
			testCase.outputPass.Logo, testCase.outputPass.Logo2x, testCase.outputPass.Strip,
			testCase.outputPass.Strip2x}

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
		"cafeid",
		"type",
		"loyaltyinfo",
		"published",
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

		row := []driver.Value{testCase.outputPass.ApplePassID, testCase.outputPass.CafeID,
			testCase.outputPass.Type, testCase.outputPass.LoyaltyInfo, testCase.outputPass.Published,
			testCase.outputPass.Design, testCase.outputPass.Icon, testCase.outputPass.Icon2x,
			testCase.outputPass.Logo, testCase.outputPass.Logo2x, testCase.outputPass.Strip,
			testCase.outputPass.Strip2x}

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

func TestGetByCafeID(t *testing.T) {
	type addTestCase struct {
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
		"cafeid",
		"type",
		"loyaltyinfo",
		"published",
		"design",
		"icon",
		"icon2x",
		"logo",
		"logo2x",
		"strip",
		"strip2x",
	}

	query := `SELECT * FROM ApplePass WHERE CafeID=$1 AND Type=$2 AND published=$3`

	testCases := []addTestCase{
		//Test OK
		{
			outputPass: outputPass,
			err:        nil,
		},
		//Test not found
		{
			outputPass: passKitModels.ApplePassDB{},
			err:        sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		row := []driver.Value{testCase.outputPass.ApplePassID, testCase.outputPass.CafeID,
			testCase.outputPass.Type, testCase.outputPass.LoyaltyInfo, testCase.outputPass.Published,
			testCase.outputPass.Design, testCase.outputPass.Icon, testCase.outputPass.Icon2x,
			testCase.outputPass.Logo, testCase.outputPass.Logo2x, testCase.outputPass.Strip,
			testCase.outputPass.Strip2x}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(row...)
			// from 1st to delete id
			mock.ExpectQuery(query).WithArgs(testCase.outputPass.CafeID, testCase.outputPass.Type,
				testCase.outputPass.Published).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(testCase.outputPass.CafeID, testCase.outputPass.Type,
				testCase.outputPass.Published).WillReturnError(testCase.err)
		}
		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		passObj, err := rep.GetPassByCafeID(context.Background(), testCase.outputPass.CafeID,
			testCase.outputPass.Type, testCase.outputPass.Published)

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
	Strip2x=NotEmpty($7, Strip2x),
    LoyaltyInfo=NotEmpty($8, LoyaltyInfo)
	WHERE CafeID=$9 AND Type=$10 AND published=$11`

	testCases := []addTestCase{
		//Test OK
		{
			inputPass: inputPass,
			err:       nil,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		args := []driver.Value{testCase.inputPass.Design, testCase.inputPass.Icon,
			testCase.inputPass.Icon2x, testCase.inputPass.Logo, testCase.inputPass.Logo2x,
			testCase.inputPass.Strip, testCase.inputPass.Strip2x, testCase.inputPass.LoyaltyInfo,
			testCase.inputPass.CafeID, testCase.inputPass.Type, testCase.inputPass.Published}

		req := mock.ExpectExec(query).WithArgs(args...)
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

//ToDo rewrite and write for getMeta
func TestUpdateMeta(t *testing.T) {
	type addTestCase struct {
		cafeID    int
		inputMeta []byte
		err       error
		finalErr  error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var cafeID int
	err = faker.FakeData(&cafeID)
	assert.NoError(t, err)

	var inputMeta []byte
	err = faker.FakeData(&inputMeta)
	assert.NoError(t, err)

	emptyMetaJson, err := json.Marshal(meta.EmptyMeta)
	if err != nil {
		t.Error(err)
	}

	queryUpdate := `UPDATE ApplePassMeta
   SET meta=$1
	WHERE CafeID=$2`

	queryInsert := `INSERT INTO ApplePassMeta(
	CafeID,
   	meta)
	VALUES ($1, $2)`

	testCases := []addTestCase{
		//Test OK
		{
			cafeID:    cafeID,
			inputMeta: inputMeta,
			err:       nil,
		},
		//Test error
		{
			cafeID:    cafeID + 1,
			inputMeta: inputMeta,
			err:       sql.ErrNoRows,
			finalErr:  nil,
		},
		////Test no cafe
		{
			cafeID:    cafeID + 1,
			inputMeta: inputMeta,
			err:       sql.ErrNoRows,
			finalErr:  sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		if testCase.err == nil {
			mock.ExpectExec(queryUpdate).WithArgs(testCase.inputMeta, testCase.cafeID).WillReturnResult(sqlmock.NewResult(0, 0))
		} else {
			mock.ExpectExec(queryUpdate).WithArgs(testCase.inputMeta, testCase.cafeID).WillReturnError(testCase.err)

			if testCase.finalErr == nil {
				mock.ExpectExec(queryInsert).WithArgs(
					testCase.cafeID, emptyMetaJson).WillReturnResult(
					sqlmock.NewResult(0, 0))
			} else {
				mock.ExpectExec(queryInsert).WithArgs(testCase.cafeID, emptyMetaJson).WillReturnError(testCase.finalErr)
			}
		}
		rep := repository.NewPostgresApplePassRepository(sqlxDB)

		err := rep.UpdateMeta(context.Background(), testCase.cafeID, testCase.inputMeta)
		assert.Equal(t, testCase.finalErr, err, message)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetMeta(t *testing.T) {
	type addTestCase struct {
		cafeID         int
		outputMeta     passKitModels.ApplePassMeta
		outputMetaJson []byte
		err            error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var cafeID int
	err = faker.FakeData(&cafeID)
	assert.NoError(t, err)

	outputMeta := passKitModels.ApplePassMeta{
		CafeID: cafeID,
		Meta: map[string]interface{}{
			"PassesCount": float64(1345),
		},
	}

	outputMetaJson, err := json.Marshal(outputMeta.Meta)
	assert.NoError(t, err)

	query := `SELECT meta FROM ApplePassMeta WHERE CafeID=$1`

	columnNames := []string{
		"meta",
	}

	testCases := []addTestCase{
		//Test OK
		{
			cafeID:         cafeID,
			outputMeta:     outputMeta,
			outputMetaJson: outputMetaJson,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		req := mock.ExpectQuery(query).WithArgs(testCase.cafeID)

		if testCase.err == nil {
			//args := []driver.Value{testCase.outputMetaJson}
			rows := sqlmock.NewRows(columnNames).AddRow(testCase.outputMetaJson)
			req.WillReturnRows(rows)
		} else {
			req.WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresApplePassRepository(sqlxDB)
		realMeta, err := rep.GetMeta(context.Background(), testCase.cafeID)
		assert.Equal(t, testCase.err, err, message)

		assert.Equal(t, testCase.outputMeta, realMeta)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
