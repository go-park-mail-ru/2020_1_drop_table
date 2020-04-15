package update_functions

import (
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResource(t *testing.T) {
	var oldVal int
	err := faker.FakeData(&oldVal)
	assert.NoError(t, err)

	newVal, err := UpdateVarPassesCount(oldVal)
	assert.NoError(t, err)

	assert.Equal(t, newVal, oldVal+1)

	var incorrectData string
	err = faker.FakeData(&incorrectData)
	assert.NoError(t, err)

	_, err = UpdateVarPassesCount(incorrectData)
	assert.Equal(t, err, ErrNotInt)

}
