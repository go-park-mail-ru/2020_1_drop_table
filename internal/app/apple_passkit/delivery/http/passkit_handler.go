package http

import (
	"2020_1_drop_table/internal/app/apple_passkit"
	"2020_1_drop_table/internal/app/apple_passkit/models"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/pkg/permissions"
	"2020_1_drop_table/internal/pkg/responses"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type applePassKitHandler struct {
	passesUsecace apple_passkit.Usecase
}

func NewPassKitHandler(r *mux.Router, us apple_passkit.Usecase) {
	handler := applePassKitHandler{
		passesUsecace: us,
	}
	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}/apple_pass",
		permissions.CheckAuthenticated(handler.UpdatePassHandler)).Methods("PUT")

	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}/apple_pass",
		permissions.CheckAuthenticated(handler.GetPassHandler)).Methods("GET")

	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}/apple_pass/{image_name}",
		permissions.CheckAuthenticated(handler.GetImageHandler)).Methods("GET")

	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}/apple_pass/new_customer",
		handler.GenerateNewPass).Methods("GET")
}

func getContent(header *multipart.FileHeader) ([]byte, error) {
	opened, err := header.Open()
	if err != nil {
		return []byte{}, nil
	}
	return ioutil.ReadAll(opened)
}

func (ap *applePassKitHandler) fetchPass(r *http.Request) (models.ApplePassDB, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return models.ApplePassDB{}, globalModels.ErrBadRequest
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" || jsonData == "null" {
		return models.ApplePassDB{}, globalModels.ErrEmptyJSON
	}

	var jsonDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &jsonDataMap); err != nil {
		return models.ApplePassDB{}, globalModels.ErrBadJSON
	}

	PassObjDB := models.ApplePassDB{Design: jsonData}

	for key, value := range r.MultipartForm.File {
		if len(value) > 1 {
			return models.ApplePassDB{}, globalModels.ErrUnexpectedFile
		}
		content, err := getContent(value[0])
		if err != nil {
			return models.ApplePassDB{}, err
		}

		switch key {
		case "icon.png":
			PassObjDB.Icon = content
		case "icon@2x.png":
			PassObjDB.Icon2x = content
		case "logo.png":
			PassObjDB.Logo = content
		case "logo@2x.png":
			PassObjDB.Logo2x = content
		case "strip.png":
			PassObjDB.Strip = content
		case "strip@2x.png":
			PassObjDB.Strip2x = content
		default:
			return models.ApplePassDB{}, fmt.Errorf(globalModels.ErrUnexpectedFilenameText, key)
		}
	}

	return PassObjDB, nil
}

func extractBoolValue(r *http.Request, valueName string) (bool, error) {
	ValueStr, ok := r.URL.Query()[valueName]
	var value bool
	var err error

	if !ok {
		value = false
	} else {
		value, err = strconv.ParseBool(ValueStr[0])
		if err != nil {
			return false, err
		} else if len(ValueStr) > 1 {
			return false, err
		}
	}

	return value, nil
}

func (ap *applePassKitHandler) UpdatePassHandler(w http.ResponseWriter, r *http.Request) {
	applePassObj, err := ap.fetchPass(r)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	publish, err := extractBoolValue(r, "publish")
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	designOnly, err := extractBoolValue(r, "design_only")
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	err = ap.passesUsecace.UpdatePass(r.Context(), applePassObj, id, publish, designOnly)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer("", w)
	return
}

func (ap *applePassKitHandler) GetPassHandler(w http.ResponseWriter, r *http.Request) {
	CafeID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	published, err := extractBoolValue(r, "published")
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	designOnly, err := extractBoolValue(r, "design_only")
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	applePassObj, err := ap.passesUsecace.GetPass(r.Context(), CafeID, published, designOnly)

	responses.SendOKAnswer(applePassObj, w)
	return
}

func (ap *applePassKitHandler) GetImageHandler(w http.ResponseWriter, r *http.Request) {
	imageName, found := mux.Vars(r)["image_name"]
	if !found {
		responses.SendSingleError(globalModels.ErrBadURLParams.Error(), w)
		return
	}

	cafeID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	published, err := extractBoolValue(r, "published")
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	image, err := ap.passesUsecace.GetImage(r.Context(), imageName, cafeID, published)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	filename := fmt.Sprintf("%s.png", imageName)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "image/png")
	http.ServeContent(w, r, filename, time.Time{}, bytes.NewReader(image))

	return
}

func (ap *applePassKitHandler) GenerateNewPass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	pass, err := ap.passesUsecace.GeneratePassObject(r.Context(), id)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	filename := "loyaltyCard.pkpass"

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/vnd.apple.pkpass")
	http.ServeContent(w, r, filename, time.Time{}, bytes.NewReader(pass.Bytes()))

	return
}
