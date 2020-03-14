package cafes

import (
	"2020_1_drop_table/mediaFiles"
	"2020_1_drop_table/owners"
	"2020_1_drop_table/projectConfig"
	"2020_1_drop_table/responses"
	"2020_1_drop_table/validators"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateCafeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.SendSingleError("bad request", w)
		return
	}

	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	owner, err := owners.StorageSession.GetOwnerByCookie(authCookie.Value)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		responses.SendSingleError("empty jsonData field", w)
		return
	}
	cafeObj := Cafe{StuffID: owner.StuffID}

	if err := json.Unmarshal([]byte(jsonData), &cafeObj); err != nil {
		responses.SendSingleError("json parsing error", w)
		return
	}

	validation, trans, err := validators.GetValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}
	if err := validation.Struct(cafeObj); err != nil {
		errs := validators.GetValidationErrors(err, trans)
		responses.SendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := mediaFiles.ReceiveFile(file, handler, "cafes")
		if err == nil {
			cafeObj.Photo = fmt.Sprintf("%s/%s", projectConfig.ServerUrl, filename)
		}
	}

	cafe, err := Storage.Append(cafeObj)
	if err != nil {
		responses.SendSingleError("cafe with this this Name already existed", w)
		return
	}

	responses.SendOKAnswer(cafe, w)
	return
}

func GetCafesListHandler(w http.ResponseWriter, r *http.Request) {
	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		responses.SendForbidden(w)
		return
	}
	owner, err := owners.StorageSession.GetOwnerByCookie(authCookie.Value)
	if err != nil {
		responses.SendForbidden(w)
		return
	}
	ownerCafes, err := Storage.getOwnerCafes(owner)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	responses.SendOKAnswer(ownerCafes, w)
}

func GetCafeHandler(w http.ResponseWriter, r *http.Request) {
	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	owner, err := owners.StorageSession.GetOwnerByCookie(authCookie.Value)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	cafe, err := Storage.Get(id)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	if !cafe.hasPermission(owner) {
		responses.SendForbidden(w)
		return
	}
	responses.SendOKAnswer(cafe, w)
}

func EditCafeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.SendSingleError("bad request", w)
		return
	}

	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	owner, err := owners.StorageSession.GetOwnerByCookie(authCookie.Value)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	cafeObj, err := Storage.Get(id)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	if !cafeObj.hasPermission(owner) {
		responses.SendForbidden(w)
		return
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		responses.SendSingleError("empty jsonData field", w)
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &cafeObj); err != nil {
		responses.SendSingleError("json parsing error", w)
		return
	}

	validation, trans, err := validators.GetValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}

	if err := validation.Struct(cafeObj); err != nil {
		errs := validators.GetValidationErrors(err, trans)
		responses.SendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := mediaFiles.ReceiveFile(file, handler, "cafes")
		if err == nil {
			cafeObj.Photo = fmt.Sprintf("%s/%s", projectConfig.ServerUrl, filename)
		}
	}

	cafeObj, err = Storage.Set(id, cafeObj)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(cafeObj, w)
}
