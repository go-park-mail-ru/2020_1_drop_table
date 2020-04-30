package loyalty_systems

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCashBack_UpdatingPass(t *testing.T) {
	type testCase struct {
		reqMap map[string]int
		DBMap  map[string]int
		result map[string]int
		err    error
	}

	c := CashBack{
		InfoVarName:   "cashback",
		PointsVarName: "points_count",
	}

	var reqInt int
	err := faker.FakeData(&reqInt)
	assert.NoError(t, err)

	var DBInt int
	err = faker.FakeData(&DBInt)
	assert.NoError(t, err)

	testCases := []testCase{
		//Test OK
		{
			reqMap: map[string]int{c.InfoVarName: reqInt},
			DBMap:  map[string]int{c.InfoVarName: DBInt},
			result: map[string]int{c.InfoVarName: reqInt},
			err:    nil,
		},
	}

	for i, test := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		var reqMapString []byte
		reqMapString, err = json.Marshal(test.reqMap)
		assert.NoError(t, err, message)

		var DBMapString []byte
		DBMapString, err = json.Marshal(test.DBMap)
		assert.NoError(t, err, message)

		result, err := c.UpdatingPass(string(reqMapString), string(DBMapString))
		assert.Equal(t, test.err, err, message)

		if test.err == nil {
			var resultMapString []byte
			resultMapString, err = json.Marshal(test.result)
			assert.NoError(t, err, message)

			assert.Equal(t, string(resultMapString), result)
		}
	}
}

func TestCashBack_CreatingCustomer(t *testing.T) {
	type testCase struct {
		loyaltyInfo    string
		customerPoints string
		err            error
	}

	c := CashBack{
		InfoVarName:   "cashback",
		PointsVarName: "points_count",
	}

	testCases := []testCase{
		//Test OK
		{
			loyaltyInfo:    "",
			customerPoints: fmt.Sprintf(`{"%s": 0}`, c.PointsVarName),
			err:            nil,
		},
	}

	for i, test := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		customerPoints, newLotaltyInfo, err := c.CreatingCustomer(test.loyaltyInfo)
		assert.Equal(t, test.err, err, message)

		if test.err == nil {
			assert.Equal(t, test.customerPoints, customerPoints)
			assert.Equal(t, test.loyaltyInfo, newLotaltyInfo)
		}
	}
}

func TestCashBack_SettingPoints(t *testing.T) {
	type testCase struct {
		reqPoints string
		newPoints string
		err       error
	}

	c := CashBack{
		InfoVarName:   "cashback",
		PointsVarName: "points_count",
	}

	testCases := []testCase{
		//Test OK
		{
			reqPoints: fmt.Sprintf(`{"%s": 10}`, c.PointsVarName),
			newPoints: fmt.Sprintf(`{"%s": 10}`, c.PointsVarName),
			err:       nil,
		},
	}

	for i, test := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		newPoints, err := c.SettingPoints("", "", test.reqPoints)
		assert.Equal(t, test.err, err, message)

		if test.err == nil {

			assert.Equal(t, test.newPoints, newPoints)
		}
	}
}
