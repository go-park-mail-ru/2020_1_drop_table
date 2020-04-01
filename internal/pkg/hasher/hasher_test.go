package hasher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashAndSalt(t *testing.T) {
	passHash := HashAndSalt(nil, "asd")
	assert.Equal(t, checkWithHash(passHash, "asd"), true)

	passHash = HashAndSalt(nil, "123321")
	assert.Equal(t, checkWithHash(passHash, "123321"), true)

	passHash = HashAndSalt(nil, "test_123Test")
	assert.Equal(t, checkWithHash(passHash, "test_123Test"), true)

}
