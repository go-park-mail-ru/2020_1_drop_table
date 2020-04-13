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
	outputCafe.PublishedApplePassID.Valid = true
	outputCafe.SavedApplePassID.Valid = true
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
		"publishedapplepassid",
		"savedapplepassid",
	}

	query := `INSERT INTO Cafe(
	CafeName, 
	Address, 
	Description, 
	StaffID, 
	OpenTime, 
	CloseTime, 
	Photo) 
	VALUES ($1,$2,$3,$4,$5,$6,$7) 
	RETURNING *`

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

		args := []driver.Value{testCase.outputCafe.CafeID, testCase.outputCafe.CafeName,
			testCase.outputCafe.Address, testCase.outputCafe.Description, testCase.outputCafe.StaffID,
			testCase.outputCafe.OpenTime, testCase.outputCafe.CloseTime, testCase.outputCafe.Photo,
			testCase.outputCafe.PublishedApplePassID, testCase.outputCafe.SavedApplePassID}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(args...)
			// from 1st to delete id
			// till the second before end to delete apple passes IDs
			mock.ExpectQuery(query).WithArgs(args[1 : len(args)-2]...).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(args[1 : len(args)-2]...).WillReturnError(testCase.err)
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
	inputCafe.PublishedApplePassID.Valid = true
	inputCafe.SavedApplePassID.Valid = true

	columnNames := []string{
		"cafeid",
		"cafename",
		"address",
		"description",
		"staffid",
		"opentime",
		"closetime",
		"photo",
		"publishedapplepassid",
		"savedapplepassid",
	}

	query := `SELECT * FROM Cafe WHERE CafeID=$1`

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
			rows := sqlmock.NewRows(columnNames).AddRow(testCase.cafe.CafeID, testCase.cafe.CafeName,
				testCase.cafe.Address, testCase.cafe.Description, testCase.cafe.StaffID, testCase.cafe.OpenTime,
				testCase.cafe.CloseTime, testCase.cafe.Photo, testCase.cafe.PublishedApplePassID,
				testCase.cafe.SavedApplePassID)

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
		inputCafeArray[i].PublishedApplePassID.Valid = true
		inputCafeArray[i].SavedApplePassID.Valid = true
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
		"publishedapplepassid",
		"savedapplepassid",
	}

	query := `SELECT * FROM Cafe WHERE StaffID=$1 ORDER BY CafeID`

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
					cafe.CloseTime, cafe.Photo, cafe.PublishedApplePassID.Int64, cafe.SavedApplePassID.Int64)
			}

			mock.ExpectQuery(query).WithArgs(testCase.staffID).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(testCase.staffID).WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresCafeRepository(sqlxDB)

		cafesObj, err := rep.GetByOwnerID(context.Background(), testCase.staffID)
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
	cafe.PublishedApplePassID.Valid = true
	cafe.SavedApplePassID.Valid = true

	columnNames := []string{
		"cafeid",
		"cafename",
		"address",
		"description",
		"staffid",
		"opentime",
		"closetime",
		"photo",
		"publishedapplepassid",
		"savedapplepassid",
	}

	query := `UPDATE Cafe SET 
	CafeName=$1, 
	Address=$2, 
	Description=$3, 
	OpenTime=$4, 
	CloseTime=$5, 
	Photo=NotEmpty($6,Photo) 
	WHERE CafeID=$7
	RETURNING *`

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
			testCase.cafe.CloseTime, testCase.cafe.Photo, testCase.cafe.PublishedApplePassID,
			testCase.cafe.SavedApplePassID}

		rows := []driver.Value{testCase.cafe.CafeID, testCase.cafe.CafeName,
			testCase.cafe.Address, testCase.cafe.Description, testCase.cafe.StaffID,
			testCase.cafe.OpenTime, testCase.cafe.CloseTime, testCase.cafe.Photo,
			testCase.cafe.PublishedApplePassID, testCase.cafe.SavedApplePassID}

		if testCase.err == nil {
			rows := sqlmock.NewRows(columnNames).AddRow(rows...)
			// append is needed to make cafeID last param
			// till the second before end to delete apple passes IDs
			mock.ExpectQuery(query).WithArgs(append(args[1:len(args)-2], args[0])...).WillReturnRows(rows)
		} else {
			mock.ExpectQuery(query).WithArgs(append(args[1:len(args)-2], args[0])...).WillReturnError(testCase.err)
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

func TestUpdateSavedPass(t *testing.T) {
	type updateSavedPassTestCase struct {
		cafe cafeModels.Cafe
		err  error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var newCafe cafeModels.Cafe
	err = faker.FakeData(&newCafe)
	assert.NoError(t, err)
	newCafe.PublishedApplePassID.Valid = true
	newCafe.SavedApplePassID.Valid = true

	query := `UPDATE Cafe SET 
    SavedApplePassID=$1
    WHERE CafeID=$2`

	testCases := []updateSavedPassTestCase{
		//Test OK
		{
			cafe: newCafe,
			err:  nil,
		},
		//Test not found
		{
			cafe: newCafe,
			err:  sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		if testCase.err == nil {
			mock.ExpectExec(query).WithArgs(testCase.cafe.SavedApplePassID,
				testCase.cafe.CafeID).WillReturnResult(sqlmock.NewResult(0, 0))
		} else {
			mock.ExpectExec(query).WithArgs(testCase.cafe.SavedApplePassID,
				testCase.cafe.CafeID).WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresCafeRepository(sqlxDB)

		err := rep.UpdateSavedPass(context.Background(), testCase.cafe)
		assert.Equal(t, testCase.err, err, message)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestUpdatePublishedPass(t *testing.T) {
	type updatePublishedPassTestCase struct {
		cafe cafeModels.Cafe
		err  error
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	var newCafe cafeModels.Cafe
	err = faker.FakeData(&newCafe)
	assert.NoError(t, err)
	newCafe.PublishedApplePassID.Valid = true
	newCafe.SavedApplePassID.Valid = true

	query := `UPDATE Cafe SET 
	PublishedApplePassID=$1
	WHERE CafeID=$2`

	testCases := []updatePublishedPassTestCase{
		//Test OK
		{
			cafe: newCafe,
			err:  nil,
		},
		//Test not found
		{
			cafe: newCafe,
			err:  sql.ErrNoRows,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		if testCase.err == nil {
			mock.ExpectExec(query).WithArgs(testCase.cafe.PublishedApplePassID,
				testCase.cafe.CafeID).WillReturnResult(sqlmock.NewResult(0, 0))
		} else {
			mock.ExpectExec(query).WithArgs(testCase.cafe.PublishedApplePassID,
				testCase.cafe.CafeID).WillReturnError(testCase.err)
		}

		rep := repository.NewPostgresCafeRepository(sqlxDB)

		err := rep.UpdatePublishedPass(context.Background(), testCase.cafe)
		assert.Equal(t, testCase.err, err, message)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
