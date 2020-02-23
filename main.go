package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	"net/http"
	"sync"
	"time"
)

//ToDo make photos available
type Owner struct {
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8,max=100"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

type OwnersStorage struct {
	sync.Mutex
	owners []Owner
}

func NewOwnersStorage() *OwnersStorage {
	return &OwnersStorage{}
}

func (ds *OwnersStorage) append(value Owner) {
	ds.owners = append(ds.owners, value)
}

func (ds *OwnersStorage) get(index int) (Owner, error) {
	if ds.count() > index {
		item := ds.owners[index]
		return item, nil
	}
	return Owner{}, nil
}

func (ds *OwnersStorage) count() int {
	return len(ds.owners)
}

func (ds *OwnersStorage) isRegistered(email string) bool {
	for i := 0; i < ds.count(); i++ {
		owner, _ := ds.Get(i)
		if owner.Email == email {
			return true
		}
	}
	return false
}

func (ds *OwnersStorage) Append(value Owner) error {
	if ds.isRegistered(value.Email) {
		err := errors.New("user with this email already existed")
		return err
	}

	ds.Lock()
	defer ds.Unlock()
	ds.append(value)

	return nil
}

func (ds *OwnersStorage) Get(index int) (Owner, error) {
	ds.Lock()
	defer ds.Unlock()
	return ds.get(index)
}

func (ds *OwnersStorage) Print() {
	ds.Lock()
	defer ds.Unlock()
	fmt.Println(ds.owners)
}

func (ds *OwnersStorage) Count() int {
	ds.Lock()
	defer ds.Unlock()
	return ds.count()
}

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Error: '%s', with status code: %d", e.Message, e.Code)
}

func createNewHttpError(code int, message string) *HttpError {
	return &HttpError{
		Code:    code,
		Message: message,
	}
}

type HttpErrorsSlice struct {
	Errors []HttpError `json:"errors"`
}

//ToDo refactor function
func getValidator() (*validator.Validate, ut.Translator, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	locale := "en"
	trans, found := uni.GetTranslator(locale)
	if !found {
		err := errors.New(fmt.Sprintf("translator for %s not found", locale))
		return nil, nil, err
	}

	v := validator.New()

	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, nil, err
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})
	return v, trans, nil
}

func getValidationErrors(err error, trans ut.Translator) []HttpError {
	errorsCount := len(err.(validator.ValidationErrors))
	errs := make([]HttpError, errorsCount, errorsCount)

	for i, e := range err.(validator.ValidationErrors) {
		validationError := createNewHttpError(400, e.Translate(trans))
		errs[i] = *validationError
	}
	return errs
}

func sendServerError(errorMessage string, w http.ResponseWriter) {
	log.Error().Msg(errorMessage)
	w.WriteHeader(500)
}

func sendSingleError(errorMessage string, w http.ResponseWriter) {
	errs := make([]HttpError, 1, 1)
	errs[0] = HttpError{400, errorMessage}
	sendSeveralErrors(errs, w)
}

func sendSeveralErrors(errors []HttpError, w http.ResponseWriter) {
	errs := HttpErrorsSlice{}
	for _, err := range errors {
		errs.Errors = append(errs.Errors, err)
	}
	serializedError, err := json.Marshal(errs)
	if err != nil {
		message := fmt.Sprintf("HttpError is json serializing: %s", err.Error())
		sendServerError(message, w)
		return
	}

	_, err = w.Write(serializedError)
	if err != nil {
		message := fmt.Sprintf("HttpError while writing is socket: %s", err.Error())
		sendServerError(message, w)
		return
	}
}

func registerView(w http.ResponseWriter, r *http.Request) {
	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		sendSingleError("empty jsonData field", w)
		return
	}
	ownerObj := Owner{CreatedAt: time.Now()}

	if err := json.Unmarshal([]byte(jsonData), &ownerObj); err != nil {
		sendSingleError("json parsing error", w)
		return
	}

	validation, trans, err := getValidator()
	if err != nil {
		message := fmt.Sprintf("HttpError in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}

	if err := validation.Struct(ownerObj); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}

	if err := owners.Append(ownerObj); err != nil {
		sendSingleError("User with this email already existed", w)
	}
	return
}

var owners = NewOwnersStorage()

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/owner", registerView).Methods("POST")
	http.Handle("/", r)
	log.Info().Msg("starting server at :8080")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Error().Msg(srv.ListenAndServe().Error())

}