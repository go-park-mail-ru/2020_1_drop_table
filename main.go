package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ====================Owner and owner storage======================

//ToDo make photos available
type Owner struct {
	ID        int       `json:"id"`
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

func (ds *OwnersStorage) append(value Owner) Owner {
	value.ID = ds.count()
	ds.owners = append(ds.owners, value)
	return value
}

func (ds *OwnersStorage) get(index int) (Owner, error) {
	if ds.count() > index && index > 0 {
		item := ds.owners[index]
		return item, nil
	}
	notFoundErrorMessage := fmt.Sprintf("Owner not fount")
	return Owner{}, errors.New(notFoundErrorMessage)
}

func (ds *OwnersStorage) count() int {
	return len(ds.owners)
}

func (ds *OwnersStorage) isRegistered(email, password string) (int, Owner) {
	for i := 0; i < ds.count(); i++ {
		owner, _ := ds.Get(i)
		if owner.Email == email && owner.Password == password {
			return 2, owner
		} else if owner.Email == email {
			return 1, Owner{}
		}
	}
	return 0, Owner{}
}

func (ds *OwnersStorage) Append(value Owner) (error, Owner) {
	if n, _ := ds.isRegistered(value.Email, ""); n > 0 {
		err := errors.New("user with this email already existed")
		return err, Owner{}
	}

	ds.Lock()
	defer ds.Unlock()
	value = ds.append(value)
	return nil, value
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

func (ds *OwnersStorage) Existed(email string, password string) (bool, Owner) {
	code, owner := ds.isRegistered(email, password)
	return code == 2, owner
}

// ====================Session and SessionStorage======================

type Session struct {
	UserID      int
	CookieToken string
}

type SessionsStorage struct {
	sync.Mutex
	sessions []Session
}

func NewSessionsStorage() *SessionsStorage {
	return &SessionsStorage{}
}

func (s *SessionsStorage) Count() int {
	s.Lock()
	defer s.Unlock()
	return len(s.sessions)
}

func (s *SessionsStorage) get(index int) (Session, error) {
	if len(s.sessions) > index {
		item := s.sessions[index]
		return item, nil
	}
	return Session{}, nil
}

func (s *SessionsStorage) Get(index int) (Session, error) {
	s.Lock()
	defer s.Unlock()
	return s.get(index)
}

func (s *SessionsStorage) createNewSession(userID int) string {
	u, _ := uuid.NewV4()
	session := Session{
		UserID:      userID,
		CookieToken: u.String(),
	}
	s.sessions = append(s.sessions, session)
	return u.String()
}

func (s *SessionsStorage) CreateNewSession(value Owner) string {
	s.Lock()
	defer s.Unlock()
	return s.createNewSession(value.ID)
}

func (s *SessionsStorage) Login(email string, password string) (string, error) {
	existed, owner := owners.Existed(email, password)
	if !existed {
		err := errors.New("user with given login and password does not exist")
		return "", err
	}
	sessionToken := s.CreateNewSession(owner)
	return sessionToken, nil
}

func (s *SessionsStorage) getOwnerByCookie(cookie string) (Owner, error) {
	for i := 0; i < s.Count(); i++ {
		session, _ := s.Get(i)
		if session.CookieToken == cookie {
			return owners.Get(session.UserID)
		}
	}
	return owners.Get(-1)
}

func hasPermission(owner Owner, cookie string) bool {
	actualOwner, err := sessions.getOwnerByCookie(cookie)
	if err != nil {
		return false
	}
	return actualOwner.ID == owner.ID
}

// ====================HttpResponses======================

type HttpResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HttpResponseError) Error() string {
	return fmt.Sprintf("Error: '%s', with status code: %d", e.Message, e.Code)
}

func createNewHttpError(code int, message string) *HttpResponseError {
	return &HttpResponseError{
		Code:    code,
		Message: message,
	}
}

type HttpErrorsSlice struct {
	Errors []HttpResponseError `json:"errors"`
}

func sendServerError(errorMessage string, w http.ResponseWriter) {
	log.Error().Msg(errorMessage)
	w.WriteHeader(http.StatusInternalServerError)
}

func sendSingleError(errorMessage string, w http.ResponseWriter) {
	log.Error().Msg(errorMessage)
	errs := make([]HttpResponseError, 1, 1)
	errs[0] = HttpResponseError{400, errorMessage}
	sendSeveralErrors(errs, w)
}

func sendSeveralErrors(errors []HttpResponseError, w http.ResponseWriter) {
	errs := HttpErrorsSlice{Errors: errors}
	serializedError, err := json.Marshal(errs)
	if err != nil {
		message := fmt.Sprintf("HttpResponseError is json serializing: %s", err.Error())
		sendServerError(message, w)
		return
	}

	_, err = w.Write(serializedError)
	if err != nil {
		message := fmt.Sprintf("HttpResponseError while writing is socket: %s", err.Error())
		sendServerError(message, w)
		return
	}
	log.Error().Msg("Validation error message sent")
}

func sendOKAnswer(data interface{}, w http.ResponseWriter) {
	type response struct {
		Data   interface{} `json:"data"`
		Errors []error     `json:"errors"`
	}
	serializedData, _ := json.Marshal(response{Data: data})
	_, err := w.Write(serializedData)
	if err != nil {
		message := fmt.Sprintf("HttpResponseError while writing is socket: %s", err.Error())
		sendServerError(message, w)
		return
	}
	log.Error().Msg("OK message sent")
}

// ====================Validator======================

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
		return ut.Add("required", "{{0} is a required field aaa}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})
	return v, trans, nil
}

func getValidationErrors(err error, trans ut.Translator) []HttpResponseError {
	errorsCount := len(err.(validator.ValidationErrors))
	errs := make([]HttpResponseError, errorsCount, errorsCount)

	for i, e := range err.(validator.ValidationErrors) {
		validationError := createNewHttpError(400, e.Translate(trans))
		errs[i] = *validationError
	}
	return errs
}

// ====================Handlers======================

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
		message := fmt.Sprintf("HttpResponseError in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}

	if err := validation.Struct(ownerObj); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}
	err, owner := owners.Append(ownerObj)
	if err != nil {
		sendSingleError("User with this email already existed", w)
		return
	}
	sendOKAnswer(owner, w)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		message := fmt.Sprintf("HttpResponseError while writing is socket: %s", err.Error())
		sendServerError(message, w)
		return
	}

	if len(data) == 0 {
		sendSingleError("no JSON body received", w)
		return
	}

	type loginForm struct {
		Email    string `validate:"required"`
		Password string `validate:"required"`
	}
	var form loginForm
	err = json.Unmarshal(data, &form)
	if err != nil {
		message := fmt.Sprintf("HttpResponseError while unmarshelling: %s", err.Error())
		sendServerError(message, w)
		return
	}

	validation, trans, err := getValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponseError in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}
	if err := validation.Struct(form); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}
	token, err := sessions.Login(form.Email, form.Password)
	if err != nil {
		sendSingleError(err.Error(), w)
		return
	}
	cookie := http.Cookie{
		Name:     "authCookie",
		Value:    token,
		Expires:  time.Time{}.AddDate(0, 1, 0),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func sendForbidden(w http.ResponseWriter) {
	sendSingleError("no permissions", w)
}

func getOwnerHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	owner, err := owners.Get(id)
	if err != nil {
		sendForbidden(w)
		return
	}
	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		sendForbidden(w)
		return
	}
	if !hasPermission(owner, authCookie.Value) {
		sendForbidden(w)
		return
	}
	sendOKAnswer(owner, w)
}

// ====================Middleware======================

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("URL: %s, METHOD: %s", r.RequestURI, r.Method)
		log.Info().Msg(msg)

		next.ServeHTTP(w, r)
	})
}

// ====================Storage======================
var owners = NewOwnersStorage()
var sessions = NewSessionsStorage()

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/owner", registerHandler).Methods("POST")
	r.HandleFunc("/api/v1/owner/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/v1/owner/{id:[0-9]+}", getOwnerHandler).Methods("GET")

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(loggingMiddleware)

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
