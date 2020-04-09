package apple_pass_generator

import (
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResource(t *testing.T) {
	envValues := make(map[string]string)
	err := faker.FakeData(&envValues)
	assert.NoError(t, err)

	var files map[string][]byte
	err = faker.FakeData(&files)
	assert.NoError(t, err)

	var file []byte
	err = faker.FakeData(&file)
	assert.NoError(t, err)

	var filename string
	err = faker.FakeData(&filename)
	assert.NoError(t, err)
	filename += ".png"

	var expectedString string
	var inputString string

	inputMap := make(map[string]interface{})
	for key, value := range envValues {
		inputString += fmt.Sprintf("<<%s>>", key)
		expectedString += fmt.Sprintf("%s", value)
		inputMap[key] = value
	}
	applePass := NewApplePass(inputString, files, inputMap)
	applePass.insertValues()
	applePass.AddImage(filename, file)
	applePass.addDesignToFiles()

	assert.Equal(t, expectedString, string(applePass.files["pass.json"]))
}

func TestGenerator(t *testing.T) {
	envValues := make(map[string]string)
	err := faker.FakeData(&envValues)
	assert.NoError(t, err)

	var files map[string][]byte
	err = faker.FakeData(&files)
	assert.NoError(t, err)

	var file []byte
	err = faker.FakeData(&file)
	assert.NoError(t, err)

	var filename string
	err = faker.FakeData(&filename)
	assert.NoError(t, err)
	filename += ".png"

	var expectedString string
	var inputString string

	inputMap := make(map[string]interface{})
	for key, value := range envValues {
		inputString += fmt.Sprintf("<<%s>>", key)
		expectedString += fmt.Sprintf("%s", value)
		inputMap[key] = value
	}
	applePassResource := NewApplePass(inputString, files, inputMap)

	generator := NewGenerator("", "", "", "")
	_, err = generator.CreateNewPass(applePassResource)
	assert.Error(t, err)

	_, err = generator.createZip(applePassResource.files)
	assert.NoError(t, err)
}
