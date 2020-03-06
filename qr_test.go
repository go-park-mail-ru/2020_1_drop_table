package main

import (
	"2020_1_drop_table/utils/qr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateErr(t *testing.T) {
	_, err := qr.generate("", 1024)
	assert.Error(t, err)

	_, err = qr.generate("dir/", -1024)
	assert.Error(t, err)
	
	_, err = qr.generate("google.com", 256)
	assert.Error(t, err)
}

func TestGenerateOk(t *testing.T) {

	_, err := qr.generate("http://yandex.com", 256)
	assert.Nil(t, err, "No errors")
	
	_, err = qr.generate("http://google.com", -256)
	assert.Nil(t, err, "No errors")
	
	_, err = qr.generate("https://yandex.ru", 1024)
	assert.Nil(t, err, "No errors")
	
	_, err = qr.generate("https://google.com/test", 128)
	assert.Nil(t, err, "No errors")
	
	_, err = qr.generate("ftp://google.com", 16)
	assert.Nil(t, err, "No errors")
	
	_, err = qr.generate("/dir/test", 256)
	assert.Nil(t, err, "No errors")
	
	_, err = qr.generate("htqttp://test.com", 1024)
	assert.Nil(t, err, "No errors")
	
	_, err = qr.generate("https://popex.", 256)
	assert.Nil(t, err, "No errors")
}



