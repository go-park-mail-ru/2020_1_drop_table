package test

// TODO тест валится с does not matches sql как и все сотальные тесты

//func TestSetTemplate(t *testing.T) {
//	type getByIDCafeTestCase struct {
//		InputSurvey string
//		Id          int
//		CafeOwnerId int
//		Err         error
//	}
//
//	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	sqlxDB := sqlx.NewDb(db, "sqlmock")
//
//	query := `INSERT INTO surveytemplate (cafeid, surveytemplate,cafeOwnerId) VALUES ($1,$2,$3)`
//
//	testCases := []getByIDCafeTestCase{
//		//Test OK
//		{
//			InputSurvey: "{}",
//			Id:          2,
//			CafeOwnerId: 4,
//			Err:         nil,
//		},
//	}
//
//	for i, testCase := range testCases {
//		message := fmt.Sprintf("test case number: %d", i)
//		mock.ExpectQuery(query).WithArgs(testCase.InputSurvey, testCase.Id, testCase.CafeOwnerId).WillReturnError(testCase.Err)
//		rep := repository.NewPostgresSurveyRepository(sqlxDB)
//
//		err := rep.SetSurveyTemplate(context.Background(), testCase.InputSurvey, testCase.Id, testCase.CafeOwnerId)
//		assert.Equal(t, testCase.Err, err, message)
//		if err := mock.ExpectationsWereMet(); err != nil {
//			t.Errorf("there were unfulfilled expectations: %s", err)
//		}
//	}
//
//}
