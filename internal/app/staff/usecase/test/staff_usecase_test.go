package test

import (
	repository2 "2020_1_drop_table/internal/app/cafe/repository"
	"2020_1_drop_table/internal/app/staff/mocks"
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/app/staff/usecase"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func getDataBaseForAdd() (*sqlx.DB, error) {
	db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, err
}

type AddTestCase struct {
	user         models.Staff
	expectedUser models.Staff
	expectedErr  error
}

type GetByIdTestCase struct {
	id           int
	expectedUser models.Staff
	expectedErr  error
}

func TestAdd(t *testing.T) {
	notNilerr := errors.New("not nil")
	testCases := []AddTestCase{
		{
			user: models.Staff{
				Email:    "kek@kek.xyz",
				Name:     "pavlik",
				Password: "123123123asd",
			},
			expectedUser: models.Staff{
				Email:    "kek@kek.xyz",
				Name:     "pavlik",
				Password: "123123123asd",
			},
			expectedErr: nil,
		},
		{
			user: models.Staff{
				Email:    "kek",
				Name:     "pavlik",
				Password: "123123123asd",
			},
			expectedUser: models.Staff{},
			expectedErr:  notNilerr,
		},
		{
			user: models.Staff{
				Email:    "kek@kek.ua",
				Name:     "pavlik",
				Password: "1",
			},
			expectedUser: models.Staff{},
			expectedErr:  notNilerr,
		},
		{
			user: models.Staff{
				Email:    "kek@kek.ua",
				Name:     "",
				Password: "1",
			},
			expectedUser: models.Staff{},
			expectedErr:  notNilerr,
		},
		{
			user: models.Staff{
				Email:    "asndjask",
				Password: "zxc",
			},
			expectedUser: models.Staff{},
			expectedErr:  notNilerr,
		},
		{
			user:         models.Staff{},
			expectedUser: models.Staff{},
			expectedErr:  notNilerr,
		},
	}
	//
	empty := models.Staff{}
	timeout := time.Second * 4
	conn, err := getDataBaseForAdd()
	if err != nil {
		fmt.Println("err while open db", err)
	}
	srepo := mocks.Repository{}
	emptyContext := context.TODO()
	cafeRepo := repository2.NewPostgresCafeRepository(conn)
	s := usecase.NewStaffUsecase(&srepo, cafeRepo, timeout)

	for _, testCase := range testCases {
		emailMatchesWithStaff := func(staff models.Staff) bool {
			assert.Equal(t, testCase.expectedUser.Email, staff.Email)
			return testCase.expectedUser.Email == staff.Email
		}

		emailMatchesWithEmail := func(email string) bool {
			assert.Equal(t, testCase.expectedUser.Email, email)
			return testCase.expectedUser.Email == email
		}
		srepo.On("Add",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(emailMatchesWithStaff)).Return(
			testCase.expectedUser, testCase.expectedErr)
		srepo.On("GetByEmail",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(emailMatchesWithEmail)).Return(
			empty, sql.ErrNoRows)
		realUser, realErr := s.Add(emptyContext, testCase.user)
		assert.Equal(t, testCase.expectedUser.Email, realUser.Email)
		if realUser.Email == "kek@kek.xyz" {
			assert.Nil(t, realErr)
		} else {
			assert.NotNil(t, realErr)
		}

	}
}

func TestGeById(t *testing.T) {
	notNilerr := errors.New("not nil")
	testCases := []GetByIdTestCase{
		{
			id: 1,
			expectedUser: models.Staff{
				Email:    "kek@kek.xyz",
				Name:     "pavlik",
				Password: "123123123asd",
			},
			expectedErr: nil,
		},
		{
			id:           -1,
			expectedUser: models.Staff{},
			expectedErr:  notNilerr,
		},
	}
	timeout := time.Second * 4
	conn, err := getDataBaseForAdd()
	if err != nil {
		fmt.Println("err while open db", err)
	}
	srepo := mocks.Repository{}
	emptyContext := context.TODO()
	cafeRepo := repository2.NewPostgresCafeRepository(conn)
	s := usecase.NewStaffUsecase(&srepo, cafeRepo, timeout)

	for _, testCase := range testCases {
		emailMatchesWithStaff := func(staff models.Staff) bool {
			assert.Equal(t, testCase.expectedUser.Email, staff.Email)
			return testCase.expectedUser.Email == staff.Email
		}

		emailMatchesWithEmail := func(email string) bool {
			assert.Equal(t, testCase.expectedUser.Email, email)
			return testCase.expectedUser.Email == email
		}
		srepo.On("Add",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(emailMatchesWithStaff)).Return(
			testCase.expectedUser, testCase.expectedErr)
		srepo.On("GetByEmail",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(emailMatchesWithEmail)).Return(
			empty, sql.ErrNoRows)
		realUser, realErr := s.Add(emptyContext, testCase.user)
		assert.Equal(t, testCase.expectedUser.Email, realUser.Email)
		if realUser.Email == "kek@kek.xyz" {
			assert.Nil(t, realErr)
		} else {
			assert.NotNil(t, realErr)
		}

	}
}
