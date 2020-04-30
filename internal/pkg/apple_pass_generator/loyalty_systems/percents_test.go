package loyalty_systems

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPercents_UpdatingPass(t *testing.T) {
	type testCase struct {
		reqMap map[string]int
		DBMap  map[int]int
		result map[int]int
		err    error
	}

	p := Percents{
		purchasesSumVarName: "purchases_sum",
		discountVarName:     "discount",
		newPurchasesVarName: "new_purchases",
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
			reqMap: map[string]int{"not int": 4},
			DBMap:  map[int]int{20: 3},
			result: map[int]int{20: 3},
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

		result, err := p.UpdatingPass(string(reqMapString), string(DBMapString))
		assert.Equal(t, test.err, err, message)

		if test.err == nil {
			var resultMapString []byte
			resultMapString, err = json.Marshal(test.result)
			assert.NoError(t, err, message)

			assert.Equal(t, string(resultMapString), result)
		}
	}
}

func TestPercents_CreatingCustomer(t *testing.T) {
	type testCase struct {
		loyaltyInfo    string
		customerPoints string
		err            error
	}

	p := Percents{
		purchasesSumVarName: "purchases_sum",
		discountVarName:     "discount",
		newPurchasesVarName: "new_purchases",
	}

	testCases := []testCase{
		//Test OK
		{
			loyaltyInfo:    "",
			customerPoints: fmt.Sprintf(`{"%s": 0, "%s": 0}`, p.purchasesSumVarName, p.discountVarName),
			err:            nil,
		},
	}

	for i, test := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		customerPoints, newLotaltyInfo, err := p.CreatingCustomer(test.loyaltyInfo)
		assert.Equal(t, test.err, err, message)

		if test.err == nil {
			assert.Equal(t, test.customerPoints, customerPoints)
			assert.Equal(t, test.loyaltyInfo, newLotaltyInfo)
		}
	}
}

func TestPercents_SettingPoints(t *testing.T) {
	type testCase struct {
		reqPoints   string
		dbPoints    string
		loyaltyInfo string
		newPoints   string
		err         error
	}

	p := Percents{
		purchasesSumVarName: "purchases_sum",
		discountVarName:     "discount",
		newPurchasesVarName: "new_purchases",
	}

	testCases := []testCase{
		//Test OK
		{
			loyaltyInfo: `{"100": 10, "5000": 20}`,
			dbPoints:    `{"discount": 10, "purchases_sum": 3000}`,
			reqPoints:   `{"new_purchases": 5000}`,
			newPoints:   `{"discount":20,"purchases_sum":8000}`,
			err:         nil,
		},
	}

	for i, test := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		newPoints, err := p.SettingPoints(test.loyaltyInfo, test.dbPoints, test.reqPoints)
		assert.Equal(t, test.err, err, message)

		if test.err == nil {
			assert.Equal(t, test.newPoints, newPoints)
		}
	}
}
