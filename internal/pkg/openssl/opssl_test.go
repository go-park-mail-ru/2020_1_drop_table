package openssl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSSl(t *testing.T) {
	res, err := Smime([]byte{1, 2, 3}, "arg")
	emptByte := []byte{}
	assert.Nil(t, err)
	assert.Equal(t, emptByte, res)
}
