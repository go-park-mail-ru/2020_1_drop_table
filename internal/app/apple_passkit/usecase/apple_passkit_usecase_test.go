package usecase_test

import (
	"2020_1_drop_table/configs"
	passKitMocks "2020_1_drop_table/internal/app/apple_passkit/mocks"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"2020_1_drop_table/internal/app/apple_passkit/usecase"
	cafeMocks "2020_1_drop_table/internal/app/cafe/mocks"
	CustomerMocks "2020_1_drop_table/internal/app/customer/mocks"
	passGeneratorMocks "2020_1_drop_table/internal/pkg/apple_pass_generator/mocks"
	"context"
	"database/sql"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestApplePassKitUsecase_GetPass(t *testing.T) {
	type GetPassTestCase struct {
		cafeID    int
		Type      string
		published bool
		passObj   models.ApplePassDB
		passMap   map[string]string
		err       error
	}

	var cafeID int
	err := faker.FakeData(&cafeID)
	assert.NoError(t, err)

	Type := "coffee_cup"

	var passObj models.ApplePassDB
	err = faker.FakeData(&passObj)
	assert.NoError(t, err)

	var published bool
	err = faker.FakeData(&published)
	assert.NoError(t, err)

	passMap := map[string]string{
		"design":       passObj.Design,
		"type":         passObj.Type,
		"loyalty_info": passObj.LoyaltyInfo,
	}
	allImages := map[string][]byte{"icon": passObj.Icon, "icon2x": passObj.Icon2x,
		"logo": passObj.Logo, "logo2x": passObj.Logo2x, "strip": passObj.Strip, "strip2x": passObj.Strip2x}
	serverStartUrl := fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s", configs.ServerUrl, configs.ApiVersion,
		cafeID, passObj.Type)
	for imageName, imageData := range allImages {
		if len(imageData) != 0 {
			passMap[imageName] = fmt.Sprintf("%s/%s", serverStartUrl, imageName)
		}
	}

	testCases := []GetPassTestCase{
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			passMap:   passMap,
		},
		//Test ErrNoRows
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   models.ApplePassDB{},
			passMap:   map[string]string{},
			err:       sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passKitMocks.Repository)
		mockCafeRepo := new(cafeMocks.Repository)
		mockPassesGenerator := new(passGeneratorMocks.Generator)
		mockPassesMeta := new(passGeneratorMocks.PassMeta)
		mockCustomerClient := new(CustomerMocks.Usecase)

		cafeIDMatches := func(id int) bool {
			assert.Equal(t, testCase.cafeID, id, message)
			return id == testCase.cafeID
		}
		passTypeMatches := func(Type string) bool {
			assert.Equal(t, testCase.Type, Type, message)
			return Type == testCase.Type
		}
		passPublishedMatches := func(published bool) bool {
			assert.Equal(t, testCase.published, published, message)
			return published == testCase.published
		}

		mockPassKitRepo.On("GetPassByCafeID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(cafeIDMatches),
			mock.MatchedBy(passTypeMatches), mock.MatchedBy(passPublishedMatches)).Return(
			testCase.passObj, testCase.err)

		passKitUsecase := usecase.NewApplePassKitUsecase(mockPassKitRepo,
			mockCafeRepo, mockCustomerClient, mockPassesGenerator, time.Second*2, mockPassesMeta)

		passMap, err := passKitUsecase.GetPass(context.Background(), testCase.cafeID,
			testCase.Type, testCase.published)

		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.passMap, passMap, message)
		}
	}
}

func TestApplePassKitUsecase_GetImage(t *testing.T) {
	type GetPassTestCase struct {
		cafeID    int
		Type      string
		published bool
		passObj   models.ApplePassDB
		imageName string
		imageData []byte
		err       error
	}

	var cafeID int
	err := faker.FakeData(&cafeID)
	assert.NoError(t, err)

	Type := "coffee_cup"

	var passObj models.ApplePassDB
	err = faker.FakeData(&passObj)
	assert.NoError(t, err)

	var published bool
	err = faker.FakeData(&published)
	assert.NoError(t, err)

	testCases := []GetPassTestCase{
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			imageName: "icon",
			imageData: passObj.Icon,
		},
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			imageName: "icon2x",
			imageData: passObj.Icon2x,
		},
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			imageName: "logo",
			imageData: passObj.Logo,
		},
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			imageName: "logo2x",
			imageData: passObj.Logo2x,
		},
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			imageName: "strip",
			imageData: passObj.Strip,
		},
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			imageName: "strip2x",
			imageData: passObj.Strip2x,
		},
		//Test ErrNoRows
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   models.ApplePassDB{},
			err:       sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passKitMocks.Repository)
		mockCafeRepo := new(cafeMocks.Repository)
		mockPassesGenerator := new(passGeneratorMocks.Generator)
		mockPassesMeta := new(passGeneratorMocks.PassMeta)
		mockCustomerClient := new(CustomerMocks.Usecase)

		cafeIDMatches := func(id int) bool {
			assert.Equal(t, testCase.cafeID, id, message)
			return id == testCase.cafeID
		}
		passTypeMatches := func(Type string) bool {
			assert.Equal(t, testCase.Type, Type, message)
			return Type == testCase.Type
		}
		passPublishedMatches := func(published bool) bool {
			assert.Equal(t, testCase.published, published, message)
			return published == testCase.published
		}

		mockPassKitRepo.On("GetPassByCafeID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(cafeIDMatches),
			mock.MatchedBy(passTypeMatches), mock.MatchedBy(passPublishedMatches)).Return(
			testCase.passObj, testCase.err)

		passKitUsecase := usecase.NewApplePassKitUsecase(mockPassKitRepo,
			mockCafeRepo, mockCustomerClient, mockPassesGenerator, time.Second*2, mockPassesMeta)

		imageData, err := passKitUsecase.GetImage(context.Background(),
			testCase.imageName, testCase.cafeID,
			testCase.Type, testCase.published)

		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.imageData, imageData, message)
		}
	}
}

func TestApplePassKitUsecase_GeneratePassObject(t *testing.T) {
	type GetPassTestCase struct {
		cafeID    int
		Type      string
		published bool
		passObj   models.ApplePassDB
		passMap   map[string]string
		err       error
	}

	var cafeID int
	err := faker.FakeData(&cafeID)
	assert.NoError(t, err)

	Type := "coffee_cup"

	var passObj models.ApplePassDB
	err = faker.FakeData(&passObj)
	assert.NoError(t, err)

	var published bool
	err = faker.FakeData(&published)
	assert.NoError(t, err)

	passMap := map[string]string{
		"design":       passObj.Design,
		"type":         passObj.Type,
		"loyalty_info": passObj.LoyaltyInfo,
	}
	allImages := map[string][]byte{"icon": passObj.Icon, "icon2x": passObj.Icon2x,
		"logo": passObj.Logo, "logo2x": passObj.Logo2x, "strip": passObj.Strip, "strip2x": passObj.Strip2x}
	serverStartUrl := fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s", configs.ServerUrl, configs.ApiVersion,
		cafeID, passObj.Type)
	for imageName, imageData := range allImages {
		if len(imageData) != 0 {
			passMap[imageName] = fmt.Sprintf("%s/%s", serverStartUrl, imageName)
		}
	}

	testCases := []GetPassTestCase{
		//Test OK
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   passObj,
			passMap:   passMap,
		},
		//Test ErrNoRows
		{
			cafeID:    cafeID,
			Type:      Type,
			published: published,
			passObj:   models.ApplePassDB{},
			passMap:   map[string]string{},
			err:       sql.ErrNoRows,
		},
	}
	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passKitMocks.Repository)
		mockCafeRepo := new(cafeMocks.Repository)
		mockPassesGenerator := new(passGeneratorMocks.Generator)
		mockPassesMeta := new(passGeneratorMocks.PassMeta)
		mockCustomerClient := new(CustomerMocks.Usecase)

		cafeIDMatches := func(id int) bool {
			assert.Equal(t, testCase.cafeID, id, message)
			return id == testCase.cafeID
		}
		passTypeMatches := func(Type string) bool {
			assert.Equal(t, testCase.Type, Type, message)
			return Type == testCase.Type
		}
		passPublishedMatches := func(published bool) bool {
			assert.Equal(t, testCase.published, published, message)
			return published == testCase.published
		}

		mockPassKitRepo.On("GetPassByCafeID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(cafeIDMatches),
			mock.MatchedBy(passTypeMatches), mock.MatchedBy(passPublishedMatches)).Return(
			testCase.passObj, testCase.err)

		passKitUsecase := usecase.NewApplePassKitUsecase(mockPassKitRepo,
			mockCafeRepo, mockCustomerClient, mockPassesGenerator, time.Second*2, mockPassesMeta)

		passMap, err := passKitUsecase.GetPass(context.Background(), testCase.cafeID,
			testCase.Type, testCase.published)

		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.passMap, passMap, message)
		}
	}
}
