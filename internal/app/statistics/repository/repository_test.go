package repository

import (
	"2020_1_drop_table/configs"
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	type getByIDCafeTestCase struct {
		cafeList []cafeModels.Cafe
		typ      string
		Since    string
		To       string
		//output   []models.StatisticsGraphRawStruct
		//err      error
		ctx context.Context
	}

	db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	cafeList := []cafeModels.Cafe{
		{
			CafeID:      1,
			CafeName:    "asd",
			Address:     "",
			Description: "",
			StaffID:     0,
			OpenTime:    time.Time{},
			CloseTime:   time.Time{},
			Photo:       "",
			Location:    "",
		},
		{
			CafeID:      2,
			CafeName:    "bsd",
			Address:     "",
			Description: "",
			StaffID:     0,
			OpenTime:    time.Time{},
			CloseTime:   time.Time{},
			Photo:       "",
			Location:    "",
		},
	}

	//columnNames := []string{
	//	"jsondata",
	//	"time",
	//	"clientuuid",
	//	"description",
	//	"staffid",
	//	"cafeid",
	//}

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	c := context.WithValue(context.Background(), configs.SessionStaffID, &session)

	testCases := []getByIDCafeTestCase{
		//Test not OK
		{
			typ:   "MONTH",
			Since: "test",
			To:    time.Now().String(),
			//output: []models.StatisticsGraphRawStruct{
			//	{
			//		Count:   1,
			//		Date:    time.Time{},
			//		CafeId:  1,
			//		StaffId: 322,
			//	},
			//	{
			//		Count:   0,
			//		Date:    time.Time{},
			//		CafeId:  2,
			//		StaffId: 228,
			//	},
			//},
			//err:      nil,
			cafeList: cafeList,
			ctx:      c,
		},
	}

	for _, testCase := range testCases {

		rep := NewPostgresStatisticsRepository(sqlxDB)
		_, err := rep.GetGraphsDataFromRepo(testCase.ctx, testCase.cafeList, testCase.typ, testCase.Since, testCase.To)
		assert.NotNil(t, err)
	}

}
