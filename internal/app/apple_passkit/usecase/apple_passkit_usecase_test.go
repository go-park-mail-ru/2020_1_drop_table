package usecase_test

import (
	"2020_1_drop_table/configs"
	passKitMocks "2020_1_drop_table/internal/app/apple_passkit/mocks"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"2020_1_drop_table/internal/app/apple_passkit/usecase"
	cafeMocks "2020_1_drop_table/internal/app/cafe/mocks"
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	CustomerMocks "2020_1_drop_table/internal/app/customer/mocks"
	CustomerModels "2020_1_drop_table/internal/app/customer/models"
	globalModels "2020_1_drop_table/internal/app/models"
	passGeneratorMocks "2020_1_drop_table/internal/pkg/apple_pass_generator/mocks"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestApplePassKitUsecase_GetPass(t *testing.T) {
	type TestCase struct {
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

	testCases := []TestCase{
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
		mockCustomerUcase := new(CustomerMocks.Usecase)

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
			mockCafeRepo, mockCustomerUcase, mockPassesGenerator, time.Second*2, mockPassesMeta)

		passMap, err := passKitUsecase.GetPass(context.Background(), testCase.cafeID,
			testCase.Type, testCase.published)

		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.passMap, passMap, message)
		}
	}
}

func TestApplePassKitUsecase_GetImage(t *testing.T) {
	type TestCase struct {
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

	testCases := []TestCase{
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
		mockCustomerUcase := new(CustomerMocks.Usecase)

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
			mockCafeRepo, mockCustomerUcase, mockPassesGenerator, time.Second*2, mockPassesMeta)

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
	type TestCase struct {
		cafeObj              cafeModels.Cafe
		passObj              models.ApplePassDB
		customerObj          CustomerModels.Customer
		metaObj              models.ApplePassMeta
		metaMap              map[string]interface{}
		Type                 string
		published            bool
		staffID              int
		bytesBuffer          *bytes.Buffer
		err                  error
		GetPassByCafeIDError error
	}

	var cafeObj cafeModels.Cafe
	err := faker.FakeData(&cafeObj)
	assert.NoError(t, err)

	Type := "coffee_cup"

	var passObj models.ApplePassDB
	err = faker.FakeData(&passObj)
	assert.NoError(t, err)
	passObj.LoyaltyInfo = `{"cups_count": 10}`
	passObj.Type = Type

	var customerObj CustomerModels.Customer
	err = faker.FakeData(&customerObj)
	assert.NoError(t, err)
	customerObj.Type = Type
	customerObj.SurveyResult = "{}"
	customerObj.CafeID = cafeObj.CafeID

	var metaObj models.ApplePassMeta
	metaObj.CafeID = cafeObj.CafeID
	metaObj.Meta = map[string]interface{}{
		"PassesCount": 8,
	}
	assert.NoError(t, err)

	passBytes, err := json.Marshal(passObj)
	assert.NoError(t, err)

	passBuffer := bytes.NewBuffer(passBytes)

	testCases := []TestCase{
		//Test OK
		{
			cafeObj:     cafeObj,
			passObj:     passObj,
			customerObj: customerObj,
			metaObj:     metaObj,
			metaMap:     map[string]interface{}{"PassesCount": 8},
			Type:        Type,
			published:   true,
			bytesBuffer: passBuffer,
			staffID:     -1,
			err:         nil,
		},
		//Test staff not published
		{
			cafeObj:     cafeObj,
			passObj:     passObj,
			customerObj: customerObj,
			metaObj:     metaObj,
			metaMap:     map[string]interface{}{"PassesCount": 8},
			Type:        Type,
			published:   false,
			bytesBuffer: passBuffer,
			staffID:     cafeObj.StaffID,
			err:         nil,
		},
		//Test not staff not published
		{
			cafeObj:     cafeObj,
			passObj:     passObj,
			customerObj: customerObj,
			metaObj:     metaObj,
			metaMap:     map[string]interface{}{"PassesCount": 8},
			Type:        Type,
			published:   false,
			bytesBuffer: passBuffer,
			staffID:     -1,
			err:         globalModels.ErrForbidden,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passKitMocks.Repository)
		mockCafeRepo := new(cafeMocks.Repository)
		mockPassesGenerator := new(passGeneratorMocks.Generator)
		mockPassesMeta := new(passGeneratorMocks.PassMeta)
		mockCustomerUcase := new(CustomerMocks.Usecase)

		cafeIDMatches := func(id int) bool {
			assert.Equal(t, testCase.cafeObj.CafeID, id, message)
			return id == testCase.cafeObj.CafeID
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
			testCase.passObj, nil)

		customerMatches := func(customer CustomerModels.Customer) bool {
			assert.Equal(t, "{}", customer.SurveyResult, message)
			assert.Equal(t, testCase.customerObj.CafeID, customer.CafeID, message)
			assert.Equal(t, testCase.customerObj.Type, customer.Type, message)
			testCase.customerObj.Points = customer.Points
			return true
		}

		mockCustomerUcase.On("Add",
			mock.AnythingOfType("*context.timerCtx"),
			mock.MatchedBy(customerMatches)).Return(
			testCase.customerObj, nil)

		mockPassKitRepo.On("GetMeta",
			mock.AnythingOfType("*context.timerCtx"),
			mock.MatchedBy(cafeIDMatches)).Return(
			testCase.metaObj, nil)

		metaMatches := func(meta map[string]interface{}) bool {
			for key := range testCase.metaMap {
				message += fmt.Sprintf(", key: %s", key)
				assert.Equal(t, testCase.metaMap[key], meta[key], message)
			}
			return true
		}

		mockPassesMeta.On("UpdateMeta",
			mock.MatchedBy(metaMatches)).Return(
			testCase.metaMap, nil)

		mockPassKitRepo.On("UpdateMeta",
			mock.AnythingOfType("*context.timerCtx"),
			mock.MatchedBy(cafeIDMatches),
			mock.AnythingOfType("[]uint8")).Return(nil)

		passKitUsecase := usecase.NewApplePassKitUsecase(mockPassKitRepo,
			mockCafeRepo, mockCustomerUcase, mockPassesGenerator, time.Second*2, mockPassesMeta)

		if !testCase.published {
			mockCafeRepo.On("GetByID",
				mock.AnythingOfType("*context.timerCtx"),
				mock.MatchedBy(cafeIDMatches)).Return(
				testCase.cafeObj, nil)
		}

		mockPassesGenerator.On("CreateNewPass",
			mock.AnythingOfType("ApplePass")).Return(
			testCase.bytesBuffer, nil)
		session := &sessions.Session{
			Values: map[interface{}]interface{}{
				"userID": testCase.staffID,
			},
		}
		ctx := context.WithValue(context.Background(), configs.SessionStaffID, session)
		passBuffer, err := passKitUsecase.GeneratePassObject(ctx, testCase.cafeObj.CafeID,
			testCase.Type, testCase.published)

		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.bytesBuffer, passBuffer, message)
		}
	}
}

func TestApplePassKitUsecase_UpdatePass(t *testing.T) {
	type TestCase struct {
		cafeObj              cafeModels.Cafe
		passObj              models.ApplePassDB
		Type                 string
		published            bool
		staffID              int
		passResponse         models.UpdateResponse
		err                  error
		GetPassByCafeIDError error
	}

	var cafeObj cafeModels.Cafe
	err := faker.FakeData(&cafeObj)
	assert.NoError(t, err)

	Type := "coffee_cup"

	var passObjPublished models.ApplePassDB
	err = faker.FakeData(&passObjPublished)
	assert.NoError(t, err)
	passObjPublished.LoyaltyInfo = `{"cups_count": 10}`
	passObjPublished.Type = Type
	passObjPublished.CafeID = cafeObj.CafeID
	passObjPublished.Published = true

	passObjNotPublished := passObjPublished
	passObjNotPublished.Published = false

	var customerObj CustomerModels.Customer
	err = faker.FakeData(&customerObj)
	assert.NoError(t, err)
	customerObj.Type = Type
	customerObj.SurveyResult = "{}"
	customerObj.CafeID = cafeObj.CafeID

	testCases := []TestCase{
		//Test OK
		{
			cafeObj:   cafeObj,
			passObj:   passObjPublished,
			Type:      Type,
			published: passObjPublished.Published,
			staffID:   cafeObj.StaffID,
			passResponse: models.UpdateResponse{
				URL: fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s/new_customer?published=true",
					configs.ServerUrl, configs.ApiVersion, passObjPublished.CafeID, passObjPublished.Type),
				QR: fmt.Sprintf("%s/media/qr/%d_%s_published.png",
					configs.ServerUrl, passObjPublished.CafeID, passObjPublished.Type),
			},
			err: nil,
		},
		//Test OK (first creation)
		{
			cafeObj:   cafeObj,
			passObj:   passObjPublished,
			Type:      Type,
			published: passObjPublished.Published,
			staffID:   cafeObj.StaffID,
			passResponse: models.UpdateResponse{
				URL: fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s/new_customer?published=true",
					configs.ServerUrl, configs.ApiVersion, passObjPublished.CafeID, passObjPublished.Type),
				QR: fmt.Sprintf("%s/media/qr/%d_%s_published.png",
					configs.ServerUrl, passObjPublished.CafeID, passObjPublished.Type),
			},
			GetPassByCafeIDError: sql.ErrNoRows,
			err:                  nil,
		},
		//Test OK (saved)
		{
			cafeObj:   cafeObj,
			passObj:   passObjNotPublished,
			Type:      Type,
			published: passObjNotPublished.Published,
			staffID:   cafeObj.StaffID,
			passResponse: models.UpdateResponse{
				URL: fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s/new_customer?published=false",
					configs.ServerUrl, configs.ApiVersion, passObjPublished.CafeID, passObjPublished.Type),
				QR: fmt.Sprintf("%s/media/qr/%d_%s_saved.png",
					configs.ServerUrl, passObjPublished.CafeID, passObjPublished.Type),
			},
			GetPassByCafeIDError: sql.ErrNoRows,
			err:                  nil,
		},
		//Test forbidden
		{
			cafeObj:   cafeObj,
			passObj:   passObjPublished,
			Type:      Type,
			published: passObjPublished.Published,
			staffID:   -1,
			passResponse: models.UpdateResponse{
				URL: fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s/new_customer?published=true",
					configs.ServerUrl, configs.ApiVersion, passObjPublished.CafeID, passObjPublished.Type),
				QR: fmt.Sprintf("%s/media/qr/%d_%s_published.png",
					configs.ServerUrl, passObjPublished.CafeID, passObjPublished.Type),
			},
			err: globalModels.ErrForbidden,
		},
	}

	for i, testCase := range testCases {
		message := fmt.Sprintf("test case number: %d", i)

		mockPassKitRepo := new(passKitMocks.Repository)
		mockCafeRepo := new(cafeMocks.Repository)
		mockPassesGenerator := new(passGeneratorMocks.Generator)
		mockPassesMeta := new(passGeneratorMocks.PassMeta)
		mockCustomerUcase := new(CustomerMocks.Usecase)

		cafeIDMatches := func(id int) bool {
			assert.Equal(t, testCase.cafeObj.CafeID, id, message)
			return id == testCase.cafeObj.CafeID
		}
		passTypeMatches := func(Type string) bool {
			assert.Equal(t, testCase.Type, Type, message)
			return Type == testCase.Type
		}
		passPublishedMatches := func(published bool) bool {
			assert.Equal(t, testCase.published, published, message)
			return published == testCase.published
		}

		mockCafeRepo.On("GetByID",
			mock.AnythingOfType("*context.timerCtx"),
			mock.MatchedBy(cafeIDMatches)).Return(
			testCase.cafeObj, nil)

		mockPassKitRepo.On("GetPassByCafeID",
			mock.AnythingOfType("*context.timerCtx"), mock.MatchedBy(cafeIDMatches),
			mock.MatchedBy(passTypeMatches), mock.MatchedBy(passPublishedMatches)).Return(
			testCase.passObj, testCase.GetPassByCafeIDError)

		if testCase.GetPassByCafeIDError != nil {
			mockPassKitRepo.On("Add",
				mock.AnythingOfType("*context.timerCtx"),
				mock.AnythingOfType("ApplePassDB")).Return(
				testCase.passObj, nil)
		}

		mockPassKitRepo.On("Update",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("ApplePassDB")).Return(nil)

		passKitUsecase := usecase.NewApplePassKitUsecase(mockPassKitRepo,
			mockCafeRepo, mockCustomerUcase, mockPassesGenerator, time.Second*2, mockPassesMeta)

		session := &sessions.Session{
			Values: map[interface{}]interface{}{
				"userID": testCase.staffID,
			},
		}
		ctx := context.WithValue(context.Background(), configs.SessionStaffID, session)
		passResponse, err := passKitUsecase.UpdatePass(ctx, testCase.passObj)

		assert.Equal(t, testCase.err, err, message)
		if err == nil {
			assert.Equal(t, testCase.passResponse, passResponse, message)
		}
	}
}
