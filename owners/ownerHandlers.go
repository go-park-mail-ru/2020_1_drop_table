package owners

import (
	"2020_1_drop_table/mediaFiles"
	"2020_1_drop_table/projectConfig"
	"2020_1_drop_table/responses"
	"2020_1_drop_table/validators"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var CookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

const CookieName = "authCookie"

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.SendSingleError("bad request", w)
		return
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" || jsonData == "null" {
		responses.SendSingleError("empty jsonData field", w)
		return
	}

	ownerObj := Owner{}

	if err := json.Unmarshal([]byte(jsonData), &ownerObj); err != nil {
		responses.SendSingleError("json parsing error", w)
		return
	}
	ownerObj.EditedAt = time.Now()

	validation, trans, err := validators.GetValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}

	if err := validation.Struct(ownerObj); err != nil {
		errs := validators.GetValidationErrors(err, trans)
		responses.SendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := mediaFiles.ReceiveFile(file, handler, "owners")
		if err == nil {
			ownerObj.Photo = fmt.Sprintf("%s/%s", projectConfig.ServerUrl, filename)
		}
	}

	owner, err := Storage.Append(ownerObj)
	if err != nil {
		responses.SendSingleError("User with this email already existed", w)
		return
	}

	session, err := CookieStore.Get(r, CookieName)
	if err != nil {
		responses.SendSingleError("bad cookies", w)
		return
	}

	session.Values["userID"] = owner.OwnerID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses.SendOKAnswer(owner, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		message := fmt.Sprintf("HttpResponse while writing is socket: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}

	if len(data) == 0 {
		responses.SendSingleError("no JSON body received", w)
		return
	}

	type loginForm struct {
		Email    string `validate:"required"`
		Password string `validate:"required"`
	}
	var form loginForm
	err = json.Unmarshal(data, &form)
	if err != nil {
		message := fmt.Sprintf("HttpResponse while unmarshelling: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}

	validation, trans, err := validators.GetValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}
	if err := validation.Struct(form); err != nil {
		errs := validators.GetValidationErrors(err, trans)
		responses.SendSeveralErrors(errs, w)
		return
	}

	existed, ownerObj, err := Storage.Existed(form.Email, form.Password)

	if !existed {
		responses.SendSingleError("no user with given login and password", w)
		return
	}

	session, err := CookieStore.Get(r, CookieName)
	if err != nil {
		responses.SendSingleError("Bad cookies", w)
		return
	}

	session.Values["userID"] = ownerObj.OwnerID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses.SendOKAnswer("", w)
}

func EditOwnerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.SendSingleError("bad request", w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	owner, err := Storage.Get(id)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	session, err := CookieStore.Get(r, CookieName)
	if err != nil {
		responses.SendSingleError("bad cookies", w)
		return
	}
	actualOwnerID, found := session.Values["userID"]
	if !found || actualOwnerID.(int) != owner.OwnerID {
		responses.SendForbidden(w)
		return
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		responses.SendSingleError("empty jsonData field", w)
		return
	}
	ownerObj := Owner{}
	if err := json.Unmarshal([]byte(jsonData), &ownerObj); err != nil {
		responses.SendSingleError("json parsing error", w)
		return
	}
	ownerObj.EditedAt = time.Now()

	validation, trans, err := validators.GetValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}

	if err := validation.Struct(ownerObj); err != nil {
		errs := validators.GetValidationErrors(err, trans)
		responses.SendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := mediaFiles.ReceiveFile(file, handler, "owners")
		if err == nil {
			ownerObj.Photo = fmt.Sprintf("%s/%s", projectConfig.ServerUrl, filename)
		}
	}

	owner, err = Storage.Set(id, ownerObj)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(owner, w)
}

func GetOwnerHandler(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, CookieName)
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

	actualOwnerID, found := session.Values["userID"]
	if !found || actualOwnerID.(int) != id {
		responses.SendForbidden(w)
		return
	}

	owner, err := Storage.Get(id)
	if err != nil {
		responses.SendForbidden(w)
		return
	}

	responses.SendOKAnswer(owner, w)
}

func GetCurrentOwnerHandler(w http.ResponseWriter, r *http.Request) {
	session, err := CookieStore.Get(r, CookieName)
	if err != nil {
		responses.SendSingleError("bad cookies", w)
		return
	}

	actualOwnerID, found := session.Values["userID"]
	if !found {
		responses.SendForbidden(w)
		return
	}

	ownerID, isInt := actualOwnerID.(int)
	if !isInt {
		responses.SendSingleError("ownerID from cookies is not int", w)
		return
	}

	owner, err := Storage.Get(ownerID)
	if err != nil {
		responses.SendForbidden(w)
		return
	}
	responses.SendOKAnswer(owner, w)
}

var Storage, _ = NewOwnerStorage("postgres", "", "5431")
