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
	"bytes"
	"context"
	"database/sql"
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
	contextTimeout  time.Duration
}

func NewApplePassKitUsecase(passKitRepo apple_passkit.Repository, cafeRepo cafe.Repository,
	customerRepo customer.Repository, passesGenerator passesGenerator.Generator,
	contextTimeout time.Duration) apple_passkit.Usecase {
	return &applePassKitUsecase{
		passKitRepo:     passKitRepo,
		cafeRepo:        cafeRepo,
		customerRepo:    customerRepo,
		passesGenerator: passesGenerator,
		contextTimeout:  contextTimeout,
	}
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

func (ap *applePassKitUsecase) addNewSavedPassToCafe(ctx context.Context, pass models.ApplePassDB,
	cafeObj cafeModels.Cafe) error {

	newPass, err := ap.passKitRepo.Add(ctx, pass)
	if err != nil {
		return err
	}

	newPassId := sql.NullInt64{
		Int64: int64(newPass.ApplePassID),
		Valid: true,
	}
	cafeObj.SavedApplePassID = newPassId
	_ = ap.cafeRepo.UpdateSavedPass(ctx, cafeObj)

	return nil
}

func (ap *applePassKitUsecase) addNewPublishedPassToCafe(ctx context.Context, pass models.ApplePassDB,
	cafeObj cafeModels.Cafe) error {

	newPass, err := ap.passKitRepo.Add(ctx, pass)
	if err != nil {
		return err
	}

	newPassId := sql.NullInt64{
		Int64: int64(newPass.ApplePassID),
		Valid: true,
	}
	cafeObj.PublishedApplePassID = newPassId
	_ = ap.cafeRepo.UpdatePublishedPass(ctx, cafeObj)

	return nil
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
		return cafeModels.Cafe{}, err
	}

	if cafeObj.StaffID != staffID {
		return cafeModels.Cafe{}, globalModels.ErrForbidden
	}

	return cafeObj, nil
}

func (ap *applePassKitUsecase) updatePass(ctx context.Context, pass models.ApplePassDB, designOnly bool) error {
	if designOnly {
		return ap.passKitRepo.UpdateDesign(ctx, pass.Design, pass.ApplePassID)
	}
	return ap.passKitRepo.Update(ctx, pass)
}

func (ap *applePassKitUsecase) UpdatePass(c context.Context, pass models.ApplePassDB, cafeID int, publish bool,
	designOnly bool) error {

	ctx, cancel := context.WithTimeout(c, ap.contextTimeout)
	defer cancel()

	cafeObj, err := ap.getOwnersCafe(ctx, cafeID)
	if err != nil {
		return err
	}

	if cafeObj.SavedApplePassID.Valid {
		pass.ApplePassID = int(cafeObj.SavedApplePassID.Int64)
		err := ap.updatePass(ctx, pass, designOnly)
		if err != nil {
			return err
		}
	} else {
		err = ap.addNewSavedPassToCafe(ctx, pass, cafeObj)
		if err != nil {
			return err
		}
	}

	if !publish {
		return nil
	}

	if cafeObj.PublishedApplePassID.Valid {
		pass.ApplePassID = int(cafeObj.PublishedApplePassID.Int64)
		err := ap.updatePass(ctx, pass, designOnly)
		if err != nil {
			return err
		}
	} else {
		err = ap.addNewPublishedPassToCafe(ctx, pass, cafeObj)
		if err != nil {
			return err
		}
	}

	return nil
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

func (ap *applePassKitUsecase) getRawPass(c context.Context, cafeID int,
	published bool) (models.ApplePassDB, error) {
	ctx, cancel := context.WithTimeout(c, ap.contextTimeout)
	defer cancel()

	cafeObj, err := ap.getOwnersCafe(ctx, cafeID)
	if err != nil {
		return models.ApplePassDB{}, err
	}

	var passID sql.NullInt64

	if published {
		passID = cafeObj.PublishedApplePassID
	} else {
		passID = cafeObj.SavedApplePassID
	}

	if !passID.Valid {
		return models.ApplePassDB{}, globalModels.ErrNoRequestedCard
	}

	return ap.passKitRepo.GetPassByID(ctx, int(passID.Int64))
}

func (ap *applePassKitUsecase) GetPass(c context.Context, cafeID int,
	published bool) (map[string]string, error) {
	passObj, err := ap.getRawPass(c, cafeID, published)
	if err != nil {
		return nil, err
	}

	return ap.getImageUrls(passObj, cafeID), nil
}

func (ap *applePassKitUsecase) GetImage(c context.Context, imageName string, cafeID int,
	published bool) ([]byte, error) {
	passObj, err := ap.getRawPass(c, cafeID, published)
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

func (ap *applePassKitUsecase) GeneratePassObject(c context.Context, cafeID int, published bool) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(c, ap.contextTimeout)
	defer cancel()

	newCustomer := customerModels.Customer{CafeID: cafeID}

	newCustomer, err := ap.customerRepo.Add(ctx, newCustomer)
	if err != nil {
		return nil, err
	}

	cafeObj, err := ap.cafeRepo.GetByID(ctx, cafeID)
	if err != nil {
		return nil, err
	}

	passMeta, err := ap.passKitRepo.UpdateMeta(ctx, cafeID)
	if err != nil {
		return nil, err
	}

	passEnv := structs.Map(passMeta)
	structs.FillMap(newCustomer, passEnv)

	cardID := -1
	if published {
		if !cafeObj.PublishedApplePassID.Valid {
			return nil, globalModels.ErrNoPublishedCard
		}
		cardID = int(cafeObj.PublishedApplePassID.Int64)
	} else {
		if !cafeObj.SavedApplePassID.Valid {
			return nil, globalModels.ErrNoPublishedCard
		}
		cardID = int(cafeObj.SavedApplePassID.Int64)

		session := ctx.Value("session").(*sessions.Session)
		staffInterface, found := session.Values["userID"]
		staffID, ok := staffInterface.(int)
		if !found || !ok || staffID != cafeObj.StaffID {
			return nil, globalModels.ErrForbidden
		}
	}

	publishedCardDB, err := ap.passKitRepo.GetPassByID(ctx, cardID)
	if err != nil {
		return nil, err
	}
	passBuffer, err := ap.passesGenerator.CreateNewPass(passDBtoPassResource(publishedCardDB, passEnv))

	return passBuffer, err
}
