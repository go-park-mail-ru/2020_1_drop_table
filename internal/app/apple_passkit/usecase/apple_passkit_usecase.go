package usecase

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app/apple_passkit"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	"2020_1_drop_table/internal/app/cafe"
	cafeModels "2020_1_drop_table/internal/app/cafe/models"
	"2020_1_drop_table/internal/app/customer"
	customerModels "2020_1_drop_table/internal/app/customer/models"
	globalModels "2020_1_drop_table/internal/app/models"
	passesGenerator "2020_1_drop_table/internal/pkg/apple_pass_generator"
	loyaltySystems "2020_1_drop_table/internal/pkg/apple_pass_generator/loyalty_systems"
	"2020_1_drop_table/internal/pkg/qr"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gorilla/sessions"
	"time"
)

type applePassKitUsecase struct {
	passKitRepo     apple_passkit.Repository
	cafeRepo        cafe.Repository
	customerRepo    customer.Repository
	passesGenerator passesGenerator.Generator
	passesMeta      passesGenerator.PassMeta
	contextTimeout  time.Duration
}

func NewApplePassKitUsecase(passKitRepo apple_passkit.Repository, cafeRepo cafe.Repository,
	customerRepo customer.Repository, passesGenerator passesGenerator.Generator,
	contextTimeout time.Duration, updateMeta passesGenerator.PassMeta) apple_passkit.Usecase {
	return &applePassKitUsecase{
		passKitRepo:     passKitRepo,
		cafeRepo:        cafeRepo,
		customerRepo:    customerRepo,
		passesGenerator: passesGenerator,
		passesMeta:      updateMeta,
		contextTimeout:  contextTimeout,
	}
}

func (ap *applePassKitUsecase) getOwnersCafe(ctx context.Context, cafeID int) (cafeModels.Cafe, error) {
	session := ctx.Value("session").(*sessions.Session)

	staffInterface, found := session.Values["userID"]
	staffID, ok := staffInterface.(int)

	if !found || !ok || staffID <= 0 {
		return cafeModels.Cafe{}, globalModels.ErrForbidden
	}

	cafeObj, err := ap.cafeRepo.GetByID(ctx, cafeID)
	if err != nil {
		fmt.Println("here")
		return cafeModels.Cafe{}, err
	}

	if cafeObj.StaffID != staffID {
		return cafeModels.Cafe{}, globalModels.ErrForbidden
	}

	return cafeObj, nil
}

func (ap *applePassKitUsecase) UpdatePass(c context.Context, pass models.ApplePassDB) (models.UpdateResponse, error) {
	ctx, cancel := context.WithTimeout(c, ap.contextTimeout)
	defer cancel()

	cafeObj, err := ap.getOwnersCafe(ctx, pass.CafeID)
	if err != nil {
		return models.UpdateResponse{}, err
	}

	if cafeObj.CafeID != pass.CafeID {
		return models.UpdateResponse{}, err
	}

	loyaltySystem, ok := loyaltySystems.LoyaltySystems[pass.Type]
	if !ok {
		return models.UpdateResponse{}, globalModels.ErrNoLoyaltyProgram
	}

	passDB, err := ap.passKitRepo.GetPassByCafeID(ctx, pass.CafeID, pass.Type, pass.Published)
	if err == sql.ErrNoRows {
		pass.LoyaltyInfo, err = loyaltySystem.UpdatingPass(pass.LoyaltyInfo, "")
		if err != nil {
			return models.UpdateResponse{}, err
		}
		passDB, err = ap.passKitRepo.Add(ctx, pass)
		if err != nil {
			return models.UpdateResponse{}, err
		}
	} else if err != nil {
		return models.UpdateResponse{}, err
	}

	pass.LoyaltyInfo, err = loyaltySystem.UpdatingPass(pass.LoyaltyInfo, passDB.LoyaltyInfo)
	if err != nil {
		return models.UpdateResponse{}, err
	}

	if pass.Published {
		pass.Published = false
		err := ap.passKitRepo.Update(ctx, pass)
		if err != nil {
			return models.UpdateResponse{}, err
		}
		pass.Published = true
	} else {
		err := ap.passKitRepo.Update(ctx, pass)
		if err != nil {
			return models.UpdateResponse{}, err
		}

		savedPassURL := fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s/new_customer?published=false",
			configs.ServerUrl, configs.ApiVersion, pass.CafeID, pass.Type)
		QrUrl := fmt.Sprintf("%s/media/qr/%d_saved.png",
			configs.ServerUrl, pass.CafeID)

		response := models.UpdateResponse{
			URL: savedPassURL,
			QR:  QrUrl,
		}
		return response, nil
	}

	err = ap.passKitRepo.Update(ctx, pass)
	if err != nil {
		return models.UpdateResponse{}, err
	}

	publishedPassURL := fmt.Sprintf("%s/%s/cafe/%d/apple_pass/%s/new_customer?published=true",
		configs.ServerUrl, configs.ApiVersion, pass.CafeID, pass.Type)
	QrUrl := fmt.Sprintf("%s/media/qr/%d_published.png",
		configs.ServerUrl, pass.CafeID)
	response := models.UpdateResponse{
		URL: publishedPassURL,
		QR:  QrUrl,
	}

	return response, nil
}

func (ap *applePassKitUsecase) getImageUrls(passObj models.ApplePassDB, cafeID int) map[string]string {
	serverStartUrl := fmt.Sprintf("%s/%s/cafe/%d/apple_pass", configs.ServerUrl, configs.ApiVersion, cafeID)
	return map[string]string{
		"design":  passObj.Design,
		"icon":    fmt.Sprintf("%s/icon", serverStartUrl),
		"icon2x":  fmt.Sprintf("%s/icon2x", serverStartUrl),
		"logo":    fmt.Sprintf("%s/logo", serverStartUrl),
		"logo2x":  fmt.Sprintf("%s/logo2x", serverStartUrl),
		"strip":   fmt.Sprintf("%s/strip", serverStartUrl),
		"strip2x": fmt.Sprintf("%s/strip2x", serverStartUrl),
	}
}

func (ap *applePassKitUsecase) GetPass(c context.Context, cafeID int, Type string,
	published bool) (map[string]string, error) {
	passObj, err := ap.passKitRepo.GetPassByCafeID(c, cafeID, Type, published)
	if err != nil {
		return nil, err
	}

	return ap.getImageUrls(passObj, cafeID), nil
}

func (ap *applePassKitUsecase) GetImage(c context.Context, imageName string, cafeID int, Type string,
	published bool) ([]byte, error) {
	passObj, err := ap.passKitRepo.GetPassByCafeID(c, cafeID, Type, published)
	if err != nil {
		return nil, err
	}
	var image []byte
	switch imageName {
	case "icon":
		image = passObj.Icon
	case "icon2x":
		image = passObj.Icon2x
	case "logo":
		image = passObj.Logo
	case "logo2x":
		image = passObj.Logo2x
	case "strip":
		image = passObj.Strip
	case "strip2x":
		image = passObj.Strip2x
	}
	if len(image) == 0 {
		return nil, globalModels.ErrNotFound
	}
	return image, nil
}

func passDBtoPassResource(db models.ApplePassDB, env map[string]interface{}) passesGenerator.ApplePass {
	files := map[string][]byte{
		"icon.png":    db.Icon,
		"icon@2x.png": db.Icon2x,
		"logo.png":    db.Logo,
		"logo@2x.png": db.Logo2x,
	}

	if len(db.Strip) != 0 && len(db.Strip2x) != 0 {
		files["strip.png"] = db.Strip
		files["strip@2x.png"] = db.Strip2x
	}

	return passesGenerator.NewApplePass(db.Design, files, env)
}

func (ap *applePassKitUsecase) updateMeta(ctx context.Context, cafeID int) (map[string]interface{}, error) {
	meta, err := ap.passKitRepo.GetMeta(ctx, cafeID)
	if err != nil {
		return nil, err
	}

	newMeta, err := ap.passesMeta.UpdateMeta(meta.Meta)
	if err != nil {
		return nil, err
	}

	jsonMeta, err := json.Marshal(newMeta)
	if err != nil {
		return nil, err
	}

	err = ap.passKitRepo.UpdateMeta(ctx, meta.CafeID, jsonMeta)
	if err != nil {
		return nil, err
	}

	return newMeta, nil
}

func (ap *applePassKitUsecase) GeneratePassObject(c context.Context, cafeID int, Type string,
	published bool) (*bytes.Buffer, error) {

	ctx, cancel := context.WithTimeout(c, ap.contextTimeout)
	defer cancel()

	loyaltySystem, ok := loyaltySystems.LoyaltySystems[Type]
	if !ok {
		return nil, globalModels.ErrNoLoyaltyProgram
	}

	publishedCardDB, err := ap.passKitRepo.GetPassByCafeID(ctx, cafeID, Type, published)
	if err != nil {
		return nil, err
	}

	customerPoints, newLoyaltyInfo, err := loyaltySystem.CreatingCustomer(publishedCardDB.LoyaltyInfo)
	if newLoyaltyInfo != publishedCardDB.LoyaltyInfo {
		publishedCardDB.LoyaltyInfo = newLoyaltyInfo
		err = ap.passKitRepo.Update(ctx, publishedCardDB)
		if err != nil {
			return nil, err
		}
	}

	newCustomer := customerModels.Customer{CafeID: cafeID, Points: customerPoints, SurveyResult: "{}"}

	newCustomer, err = ap.customerRepo.Add(ctx, newCustomer)
	if err != nil {
		return nil, err
	}

	cafeObj, err := ap.cafeRepo.GetByID(ctx, cafeID)
	if err != nil {
		return nil, err
	}

	passEnv, err := ap.updateMeta(ctx, cafeID)
	if err != nil {
		return nil, err
	}

	structs.FillMap(newCustomer, passEnv)

	if !published {
		session := ctx.Value("session").(*sessions.Session)
		staffInterface, found := session.Values["userID"]
		staffID, ok := staffInterface.(int)
		if !found || !ok || staffID != cafeObj.StaffID {
			return nil, globalModels.ErrForbidden
		}
	}

	passBuffer, err := ap.passesGenerator.CreateNewPass(passDBtoPassResource(publishedCardDB, passEnv))

	return passBuffer, err
}

func (ap *applePassKitUsecase) createQRs(cafeID int) error {
	savedPassURL := fmt.Sprintf("%s/%s/cafe/%d/apple_pass/new_customer?published=false",
		configs.ServerUrl, configs.ApiVersion, cafeID)
	savedPassPath := fmt.Sprintf("%d_saved", cafeID)
	publishedPassURL := fmt.Sprintf("%s/%s/cafe/%d/apple_pass/new_customer?published=true",
		configs.ServerUrl, configs.ApiVersion, cafeID)
	publishedPassPath := fmt.Sprintf("%d_published", cafeID)

	_, err := qr.GenerateToFile(savedPassURL, savedPassPath)
	if err != nil {
		return err
	}

	_, err = qr.GenerateToFile(publishedPassURL, publishedPassPath)
	if err != nil {
		return err
	}
	return nil
}
