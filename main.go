package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	"log"
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

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func createNewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

type ResponseErrorText struct {
	Errors []Error `json:"errors"`
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

func getValidationErrors(err error, trans ut.Translator) ResponseErrorText {
	errs := ResponseErrorText{}
	for _, e := range err.(validator.ValidationErrors) {
		validationError := createNewError(400, e.Translate(trans))

		errs.Errors = append(errs.Errors, *validationError)
	}
	return errs
}

func registerView(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		jsonData := r.FormValue("jsonData")
		if jsonData == "" {
			errs := ResponseErrorText{}
			errs.Errors = append(errs.Errors, Error{400, "empty jsonData field"})
			serializedError, err := json.Marshal(errs)
			if err != nil {
				w.WriteHeader(500)
				return
			}

			_, err = w.Write(serializedError)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			return
		}
		ownerObj := Owner{CreatedAt: time.Now()}

		if err := json.Unmarshal([]byte(jsonData), &ownerObj); err != nil {
			errs := ResponseErrorText{}
			errs.Errors = append(errs.Errors, Error{400, "json parsing error"})
			serialized, _ := json.Marshal(errs)
			_, err = w.Write(serialized)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			fmt.Println("Unmarshal error: ", err)
			return
		}

		validation, trans, validatorError := getValidator()

		if validatorError != nil {
			w.WriteHeader(500)
			return
		}

		if err := validation.Struct(ownerObj); err != nil {
			errs := getValidationErrors(err, trans)
			serialized, err := json.Marshal(errs)
			if err != nil {
				w.WriteHeader(500)
				return
			}

			_, err = w.Write(serialized)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			return
		}
		err := owners.Append(ownerObj)
		if err != nil {
			errs := ResponseErrorText{}
			errs.Errors = append(errs.Errors, Error{400, "User with this email already existed"})
			serialized, _ := json.Marshal(errs)
			_, err = w.Write(serialized)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}
		owners.Print()

	default:
		methodErr := ResponseErrorText{}
		methodErr.Errors = append(methodErr.Errors, Error{400, "this method is unavailable"})
		serialized, _ := json.Marshal(methodErr)
		_, err := w.Write(serialized)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		return
	}
	return

}

var owners = NewOwnersStorage()

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/owner", registerView)
	http.Handle("/", r)
	fmt.Println("starting server at :8080")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
