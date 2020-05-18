package repository_test

import (
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	"2020_1_drop_table/internal/app/cafe/repository"
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
	type getByIDCafeTestCase struct {
		inputCafe  cafeModels.Cafe
		outputCafe cafeModels.Cafe
		err        error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var outputCafe cafeModels.Cafe
	err = faker.FakeData(&outputCafe)
	assert.NoError(t, err)
	inputCafe := outputCafe
	inputCafe.CafeID = 0

	columnNames := []string{
		"cafeid",
		"cafename",
		"address",
		"description",
		"staffid",
		"opentime",
		"closetime",
		"photo",
		"location_str",
	}

	query := `INSERT INTO Cafe(
	CafeName,
	Address,
	Description,
	StaffID,
	OpenTime,
	CloseTime,
	Photo,
   location,
   location_str)
	VALUES ($1,$2,$3,$4,$5,$6,$7,ST_GeomFromEWKT($8),$9)
	RETURNING CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str`

	testCases := []getByIDCafeTestCase{
		//Test OK
		{
			inputCafe:  inputCafe,
			outputCafe: outputCafe,
			err:        nil,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)
		postGisPoint := repository.GeneratePointToGeoWithPoint(testCase.outputCafe.Location)
		args := []driver.Value{testCase.outputCafe.CafeName, testCase.outputCafe.Address,
			testCase.outputCafe.Description, testCase.outputCafe.StaffID,
			testCase.outputCafe.OpenTime, testCase.outputCafe.CloseTime,
			testCase.outputCafe.Photo, postGisPoint, testCase.outputCafe.Location}

		rows := []driver.Value{testCase.outputCafe.CafeID, testCase.outputCafe.CafeName,
			testCase.outputCafe.Address, testCase.outputCafe.Description, testCase.outputCafe.StaffID,
			testCase.outputCafe.OpenTime, testCase.outputCafe.CloseTime, testCase.outputCafe.Photo,
			testCase.outputCafe.Location}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(rows...)
			// from 1st to delete id
			// till the second before end to delete apple passes IDs
			mock.ExpectQuery(query).WithArgs(args...).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(args...).WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresCafeRepository(sqlxDB)

		cafeObj, err := rep.Add(context.Background(), testCase.inputCafe)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.outputCafe, cafeObj, message)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}

}

func TestGetByID(t *testing.T) {
	type getByIDCafeTestCase struct {
		cafe cafeModels.Cafe
		err  error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var inputCafe cafeModels.Cafe
	err = faker.FakeData(&inputCafe)
	assert.NoError(t, err)

	columnNames := []string{
		"cafeid",
		"cafename",
		"address",
		"description",
		"staffid",
		"opentime",
		"closetime",
		"photo",
		"location_str",
	}

	query := `SELECT CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str FROM Cafe WHERE CafeID=$1`

	testCases := []getByIDCafeTestCase{
		//Test OK
		{
			cafe: inputCafe,
			err:  nil,
		},
		//Test not found
		{
			cafe: inputCafe,
			err:  sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		if testCase.err == nil {
			args := []driver.Value{testCase.cafe.CafeID, testCase.cafe.CafeName,
				testCase.cafe.Address, testCase.cafe.Description, testCase.cafe.StaffID, testCase.cafe.OpenTime,
				testCase.cafe.CloseTime, testCase.cafe.Photo, testCase.cafe.Location}

			rows := sqlmock.NewRows(columnNames).AddRow(args...)
			mock.ExpectQuery(query).WithArgs(testCase.cafe.CafeID).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(testCase.cafe.CafeID).WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresCafeRepository(sqlxDB)

		cafeObj, err := rep.GetByID(context.Background(), testCase.cafe.CafeID)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.cafe, cafeObj, message)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByOwnerID(t *testing.T) {
	type getByOwnerIDTestCase struct {
		cafesArray []cafeModels.Cafe
		staffID    int
		err        error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var staffID int
	err = faker.FakeData(&staffID)
	assert.NoError(t, err)

	inputCafeArray := make([]cafeModels.Cafe, 5, 5)
	err = faker.FakeData(&inputCafeArray)
	assert.NoError(t, err)

	for i := range inputCafeArray {
		inputCafeArray[i].StaffID = staffID
		inputCafeArray[i].CafeID = i + 1
	}

	columnNames := []string{
		"cafeid",
		"cafename",
		"address",
		"description",
		"staffid",
		"opentime",
		"closetime",
		"photo",
		"location_str",
	}

	query := `SELECT CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str FROM Cafe WHERE StaffID=$1 ORDER BY CafeID`

	testCases := []getByOwnerIDTestCase{
		//Test OK
		{
			cafesArray: inputCafeArray,
			staffID:    staffID,
			err:        nil,
		},
		//Test not found
		{
			cafesArray: nil,
			staffID:    0,
			err:        sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames)
			for _, cafe := range testCase.cafesArray {
				rows.AddRow(cafe.CafeID, cafe.CafeName, cafe.Address, cafe.Description, cafe.StaffID, cafe.OpenTime,
					cafe.CloseTime, cafe.Photo, cafe.Location)
			}

			mock.ExpectQuery(query).WithArgs(testCase.staffID).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(testCase.staffID).WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresCafeRepository(sqlxDB)

		cafesObj, err := rep.GetByOwnerId(context.Background(), testCase.staffID)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.cafesArray, cafesObj, message)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestUpdate(t *testing.T) {
	type updateTestCase struct {
		cafe cafeModels.Cafe
		err  error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var cafe cafeModels.Cafe
	err = faker.FakeData(&cafe)
	assert.NoError(t, err)

	columnNames := []string{
		"cafeid",
		"cafename",
		"address",
		"description",
		"staffid",
		"opentime",
		"closetime",
		"photo",
		"location_str",
	}

	query := `UPDATE Cafe SET 
	CafeName=$1, 
	Address=$2, 
	Description=$3, 
	OpenTime=$4, 
	CloseTime=$5, 
	Photo=NotEmpty($6,Photo) 
	WHERE CafeID=$7
	RETURNING CafeID,CafeName,Address,Description,StaffID,OpenTime,CloseTime,Photo,location_str`

	testCases := []updateTestCase{
		//Test OK
		{
			cafe: cafe,
			err:  nil,
		},
		//Test not found
		{
			cafe: cafeModels.Cafe{},
			err:  sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		args := []driver.Value{testCase.cafe.CafeID, testCase.cafe.CafeName,
			testCase.cafe.Address, testCase.cafe.Description, testCase.cafe.OpenTime,
			testCase.cafe.CloseTime, testCase.cafe.Photo}

		rows := []driver.Value{testCase.cafe.CafeID, testCase.cafe.CafeName,
			testCase.cafe.Address, testCase.cafe.Description, testCase.cafe.StaffID,
			testCase.cafe.OpenTime, testCase.cafe.CloseTime, testCase.cafe.Photo,
			testCase.cafe.Location}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(rows...)
			// append is needed to make cafeID last param
			// till the second before end to delete apple passes IDs
			mock.ExpectQuery(query).WithArgs(append(args[1:], args[0])...).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(append(args[1:], args[0])...).WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresCafeRepository(sqlxDB)

		cafeObj, err := rep.Update(context.Background(), testCase.cafe)
		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.cafe, cafeObj, message)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
