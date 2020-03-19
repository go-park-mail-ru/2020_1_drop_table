package qr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateErr(t *testing.T) {
	_, err := Generate("", 1024)
	assert.Error(t, err)

	_, err = Generate("dir/", -1024)
	assert.Error(t, err)

	_, err = Generate("google.com", 256)
	assert.Error(t, err)
}

func TestGenerateOk(t *testing.T) {

	_, err := Generate("http://yandex.com", 256)
	assert.Nil(t, err, "No errors")

	_, err = Generate("http://google.com", -256)
	assert.Nil(t, err, "No errors")

	_, err = Generate("https://yandex.ru", 1024)
	assert.Nil(t, err, "No errors")

	_, err = Generate("https://google.com/test", 128)
	assert.Nil(t, err, "No errors")

	_, err = Generate("ftp://google.com", 16)
	assert.Nil(t, err, "No errors")

	_, err = Generate("/dir/test", 256)
	assert.Nil(t, err, "No errors")

	_, err = Generate("htqttp://test.com", 1024)
	assert.Nil(t, err, "No errors")

	_, err = Generate("https://popex.", 256)
	assert.Nil(t, err, "No errors")
}
