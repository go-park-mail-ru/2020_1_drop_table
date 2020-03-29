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

	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}/apple_pass/new_customer", handler.GenerateNewPass).Methods("GET")
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

	publishRaw, ok := r.URL.Query()["publish"]
	var publish bool
	if !ok {
		publish = false
	} else {
		publish, err = strconv.ParseBool(publishRaw[0])
		if err != nil {
			responses.SendSingleError(err.Error(), w)
			return
		} else if len(publishRaw) > 1 {
			fmt.Println(len(publishRaw))
			responses.SendSingleError(globalModels.ErrBadURLParams.Error(), w)
			return
		}
	}

	err = ap.passesUsecace.UpdatePass(r.Context(), applePassObj, id, publish, false)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer("", w)
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

}
