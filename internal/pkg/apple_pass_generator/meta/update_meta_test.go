package meta

import (
	"2020_1_drop_table/internal/pkg/apple_pass_generator/meta/update_functions"
	"errors"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMeta(t *testing.T) {
	type updateMetaTestCase struct {
		oldValues map[string]interface{}
		newValues map[string]interface{}
		err       error
	}

	var passesCount int
	err := faker.FakeData(&passesCount)
	assert.NoError(t, err)

	testCases := []updateMetaTestCase{
		// Test OK
		{
			oldValues: map[string]interface{}{"PassesCount": passesCount},
			newValues: map[string]interface{}{"PassesCount": passesCount + 1},
			err:       nil,
		},
		// Test not int
		{
			oldValues: map[string]interface{}{"PassesCount": "NOT INT"},
			newValues: nil,
			err:       update_functions.ErrNotInt,
		},
		// Test unresolved value name
		{
			oldValues: map[string]interface{}{"NO GIVEN NAME": "NOT INT"},
			newValues: nil,
			err:       errors.New(fmt.Sprintf("not found update func for var <<%s>>", "NO GIVEN NAME")),
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		newValues, err := UpdateMeta(testCase.oldValues)
		assert.Equal(t, err, testCase.err, message)

		if testCase.err == nil {
			assert.Equal(t, newValues, testCase.newValues, message)
		}
	}
}
