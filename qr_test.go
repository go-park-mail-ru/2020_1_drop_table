package main

import (
	"2020_1_drop_table/utils/qr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateErr(t *testing.T) {
	_, err := qr.Generate("", 1024)
	assert.Error(t, err)

	_, err = qr.Generate("dir/", -1024)
	assert.Error(t, err)

	_, err = qr.Generate("google.com", 256)
	assert.Error(t, err)
}

func TestGenerateOk(t *testing.T) {

	_, err := qr.Generate("http://yandex.com", 256)
	assert.Nil(t, err, "No errors")

	_, err = qr.Generate("http://google.com", -256)
	assert.Nil(t, err, "No errors")

	_, err = qr.Generate("https://yandex.ru", 1024)
	assert.Nil(t, err, "No errors")

	_, err = qr.Generate("https://google.com/test", 128)
	assert.Nil(t, err, "No errors")

	_, err = qr.Generate("ftp://google.com", 16)
	assert.Nil(t, err, "No errors")

	_, err = qr.Generate("/dir/test", 256)
	assert.Nil(t, err, "No errors")

	_, err = qr.Generate("htqttp://test.com", 1024)
	assert.Nil(t, err, "No errors")

	_, err = qr.Generate("https://popex.", 256)
	assert.Nil(t, err, "No errors")
}
