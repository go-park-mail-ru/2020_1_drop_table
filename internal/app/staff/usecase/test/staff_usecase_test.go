package test

import (
	cafeMock "2020_1_drop_table/internal/app/cafe/mocks"
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/staff/mocks"
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/app/staff/usecase"
	"2020_1_drop_table/internal/pkg/hasher"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

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

type GetByEmailAndPasswordTestCase struct {
	form         models.LoginForm
	expectedUser models.Staff
	expectedErr  error
}

type UpdateTestCase struct {
	user         models.SafeStaff
	expectedUser models.SafeStaff
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

	srepo := new(mocks.Repository)
	emptyContext := context.TODO()
	cafeRepo := cafeMock.Repository{}
	s := usecase.NewStaffUsecase(srepo, &cafeRepo, timeout)

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
			expectedErr:  nil,
		},
	}
	timeout := time.Second * 4

	srepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	s := usecase.NewStaffUsecase(&srepo, &cafeRepo, timeout)

	for _, testCase := range testCases {
		srepo.On("GetByID",
			mock.AnythingOfType("*context.timerCtx"), testCase.id).Return(
			testCase.expectedUser, testCase.expectedErr)
		sessionUserID := testCase.id
		session := sessions.Session{Values: map[interface{}]interface{}{"userID": sessionUserID}}
		c := context.WithValue(context.Background(), "session", &session)
		realUser, realErr := s.GetByID(c, testCase.id)
		assert.Equal(t, testCase.expectedUser.Email, realUser.Email)
		assert.Equal(t, testCase.expectedErr, realErr)

	}
}

func TestUpdate(t *testing.T) {
	notNilerr := errors.New("not nil")
	testCases := []UpdateTestCase{
		{ // ok 1
			user: models.SafeStaff{
				Email: "email@email.ru",
				Name:  "keklolarbidol",
			},
			expectedUser: models.SafeStaff{
				Email: "email@email.ru",
				Name:  "keklolarbidol",
			},
			expectedErr: nil,
		},

		{ //not ok 2
			user: models.SafeStaff{
				Email: "kek",
				Name:  "pavlik",
			},
			expectedUser: models.SafeStaff{},
			expectedErr:  notNilerr,
		},
	}
	//
	timeout := time.Second * 4

	srepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	s := usecase.NewStaffUsecase(&srepo, &cafeRepo, timeout)

	for _, testCase := range testCases {

		srepo.On("Update",
			mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.SafeStaff")).Return(
			testCase.expectedUser, testCase.expectedErr)

		sessionUserID := testCase.user.StaffID
		session := sessions.Session{Values: map[interface{}]interface{}{"userID": sessionUserID}}
		c := context.WithValue(context.Background(), "session", &session)
		realUser, realErr := s.Update(c, testCase.user)
		fmt.Println(realUser)
		assert.Equal(t, testCase.expectedUser.Email, realUser.Email)
		if testCase.expectedUser.Email == "email@email.ru" {
			assert.Nil(t, realErr)
		} else {
			assert.NotNil(t, realErr)
		}

	}
}

func TestGet(t *testing.T) {

	testCases := []GetByEmailAndPasswordTestCase{
		{
			form: models.LoginForm{Email: "asd@asd.ru",
				Password: "password1",
			},
			expectedUser: models.Staff{
				Email:    "",
				Name:     "pavlik",
				Password: "password1",
			},
			expectedErr: globalModels.ErrNotFound,
		},
		{
			form: models.LoginForm{Email: "bad",
				Password: "password1",
			},
			expectedUser: models.Staff{
				StaffID:  0,
				Name:     "",
				Email:    "",
				Password: "",
				EditedAt: time.Time{},
				Photo:    "",
				IsOwner:  false,
				CafeId:   0,
				Position: "",
			},
			expectedErr: globalModels.ErrNotFound,
		},
	}
	timeout := time.Second * 4
	srepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	s := usecase.NewStaffUsecase(&srepo, &cafeRepo, timeout)

	for _, testCase := range testCases {
		testCase.form.Password, _ = hasher.HashAndSalt(nil, testCase.form.Password)
		testCase.expectedUser.Password = testCase.form.Password
		srepo.On("GetByEmail",
			mock.AnythingOfType("*context.timerCtx"), testCase.form.Email).Return(
			testCase.expectedUser, testCase.expectedErr)

		realUser, realErr := s.GetByEmailAndPassword(context.TODO(), testCase.form)
		assert.Equal(t, testCase.expectedUser.Email, realUser.Email)
		assert.Equal(t, testCase.expectedErr, realErr)

	}

}

func TestGenerateQr(t *testing.T) {
	user := models.Staff{
		StaffID:  0,
		Name:     "test",
		Email:    "email@email.ru",
		Password: "",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  true,
		CafeId:   2,
		Position: "",
	}
	timeout := time.Second * 4
	srepo := mocks.Repository{}
	srepo.On("GetByID",
		mock.AnythingOfType("*context.cancelCtx"), 228).Return(
		user, nil)
	cafeRepo := cafeMock.Repository{}
	returnedCafe := cafeModels.Cafe{
		CafeID:               2,
		CafeName:             "",
		Address:              "",
		Description:          "",
		StaffID:              0,
		OpenTime:             time.Time{},
		CloseTime:            time.Time{},
		Photo:                "",
		PublishedApplePassID: sql.NullInt64{},
		SavedApplePassID:     sql.NullInt64{},
	}

	cafeRepo.On("GetByID",
		mock.AnythingOfType("*context.timerCtx"), 2).Return(
		returnedCafe, nil)
	srepo.On("AddUuid", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string"), 2).Return(nil)
	s := usecase.NewStaffUsecase(&srepo, &cafeRepo, timeout)
	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	c := context.WithValue(context.Background(), "session", &session)
	_, err := s.GetQrForStaff(c, 2)
	assert.Nil(t, err)

}

func TestDelQr(t *testing.T) {
	timeout := time.Second * 4
	srepo := mocks.Repository{}
	cafeRepo := cafeMock.Repository{}
	s := usecase.NewStaffUsecase(&srepo, &cafeRepo, timeout)
	err := s.DeleteQrCodes("not exist")
	assert.Equal(t, err.Error(), "remove media/qr/not exist.png: no such file or directory")
}

func TestGetStaffList(t *testing.T) {
	timeout := time.Second * 4
	srepo := mocks.Repository{}
	user := models.Staff{
		StaffID:  2,
		Name:     "",
		Email:    "",
		Password: "",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  true,
		CafeId:   0,
		Position: "",
	}
	srepo.On("GetByID",
		mock.AnythingOfType("*context.cancelCtx"), 2).Return(
		user, nil)
	resMap := make(map[string][]models.StaffByOwnerResponse)

	srepo.On("GetStaffListByOwnerId", mock.AnythingOfType("*context.valueCtx"), 2).Return(resMap, nil)
	cafeRepo := cafeMock.Repository{}
	s := usecase.NewStaffUsecase(&srepo, &cafeRepo, timeout)
	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 2}}
	c := context.WithValue(context.Background(), "session", &session)
	res, err := s.GetStaffListByOwnerId(c, 2)
	assert.Equal(t, resMap, res)
	assert.Nil(t, err)
}