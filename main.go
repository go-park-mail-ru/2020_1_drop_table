package main

import (
	"crypto/md5"
	"encoding/hex"
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
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//=====================Hasher func======================

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//=====================Owner and owner storage======================

//ToDo make photos available
type Owner struct {
	ID       int       `json:"id"`
	Name     string    `json:"name" validate:"required,min=4,max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	EditedAt time.Time `json:"editedAt" validate:"required"`
	Photo    string    `json:"photo"`
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

func (ds *OwnersStorage) set(i int, value Owner) Owner {
	ds.owners[i] = value
	return value
}

func (ds *OwnersStorage) get(index int) (Owner, error) {
	if ds.count() > index && index >= 0 {
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
	password = GetMD5Hash(password)
	for i := 0; i < ds.count(); i++ {
		owner, _ := ds.Get(i)
		if owner.Email == email && owner.Password == password {
			return 2, owner
		} else if owner.Email == email {
			return 1, Owner{}
		}
	}
	return -1, Owner{}
}

func (ds *OwnersStorage) Append(value Owner) (error, Owner) {
	if n, _ := ds.isRegistered(value.Email, ""); n != -1 {
		err := errors.New("user with this email already existed")
		return err, Owner{}
	}
	value.Password = GetMD5Hash(value.Password)
	ds.Lock()
	defer ds.Unlock()
	value = ds.append(value)
	return nil, value
}

func (ds *OwnersStorage) Set(i int, value Owner) (Owner, error) {
	if i > ds.Count() {
		err := errors.New(fmt.Sprintf("no user with id: %d", i))
		return Owner{}, err
	}
	value.ID = i

	ds.Lock()
	defer ds.Unlock()
	value = ds.set(i, value)
	return value, nil
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

//=====================Session and SessionStorage======================

type Session struct {
	UserID      int
	CookieToken string
	ExpiresDate time.Time
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

func (s *SessionsStorage) get(index int) Session {
	if len(s.sessions) > index {
		item := s.sessions[index]
		return item
	}
	return Session{}
}

func (s *SessionsStorage) Get(index int) Session {
	s.Lock()
	defer s.Unlock()
	return s.get(index)
}

func (s *SessionsStorage) createNewSession(userID int, expiresDate time.Time) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	session := Session{
		UserID:      userID,
		CookieToken: u.String(),
		ExpiresDate: expiresDate,
	}
	s.sessions = append(s.sessions, session)
	return u.String(), nil
}

func (s *SessionsStorage) CreateNewSession(value Owner, expiresDate time.Time) (string, error) {
	s.Lock()
	defer s.Unlock()
	return s.createNewSession(value.ID, expiresDate)
}

func (s *SessionsStorage) Login(email string, password string, expiresDate time.Time) (string, error) {
	existed, owner := owners.Existed(email, password)
	if !existed {
		err := errors.New("user with given login and password does not exist")
		return "", err
	}
	sessionToken, err := s.CreateNewSession(owner, expiresDate)
	return sessionToken, err
}

func (s *SessionsStorage) getOwnerByCookie(cookie string) (Owner, error) {
	for i := 0; i < s.Count(); i++ {
		session := s.Get(i)
		timeDiff := session.ExpiresDate.Sub(time.Now())
		if session.CookieToken == cookie && timeDiff > 0 {
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

//=====================Cafe and CafeStorage======================

//ToDo make photos available
type Cafe struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Address     string    `json:"address" validate:"required"`
	Description string    `json:"description" validate:"required"`
	OwnerID     int       `json:"ownerID"`
	OpenTime    time.Time `json:"openTime"`
	CloseTime   time.Time `json:"closeTime"`
	Photo       string    `json:"photo"`
}

func (c *Cafe) hasPermission(owner Owner) bool {
	return c.OwnerID == owner.ID
}

type CafesStorage struct {
	sync.Mutex
	cafes []Cafe
}

func NewCafesStorage() *CafesStorage {
	return &CafesStorage{}
}

func (cs *CafesStorage) append(value Cafe) Cafe {
	value.ID = cs.count()
	cs.cafes = append(cs.cafes, value)
	return value
}

func (cs *CafesStorage) set(i int, value Cafe) Cafe {
	cs.cafes[i] = value
	return value
}

func (cs *CafesStorage) get(index int) (Cafe, error) {
	if cs.count() > index && index >= 0 {
		item := cs.cafes[index]
		return item, nil
	}
	notFoundErrorMessage := fmt.Sprintf("Cafe not fount")
	return Cafe{}, errors.New(notFoundErrorMessage)
}

func (cs *CafesStorage) count() int {
	return len(cs.cafes)
}

func (cs *CafesStorage) Print() {
	cs.Lock()
	defer cs.Unlock()
	fmt.Println(cs.cafes)
}

func (cs *CafesStorage) Append(value Cafe) (error, Cafe) {
	cs.Lock()
	defer cs.Unlock()
	value = cs.append(value)
	return nil, value
}

func (cs *CafesStorage) Count() int {
	cs.Lock()
	defer cs.Unlock()
	return cs.count()
}

func (cs *CafesStorage) Get(index int) (Cafe, error) {
	cs.Lock()
	defer cs.Unlock()
	return cs.get(index)
}

func (cs *CafesStorage) getOwnerCafes(owner Owner) []Cafe {
	var ownerCafes []Cafe
	for i := 0; i < cs.Count(); i++ {
		cafe, _ := cs.Get(i)
		if cafe.OwnerID == owner.ID {
			ownerCafes = append(ownerCafes, cafe)
		}
	}
	return ownerCafes
}

func (cs *CafesStorage) Set(i int, value Cafe) (Cafe, error) {
	if i > cs.Count() {
		err := errors.New(fmt.Sprintf("no user with id: %d", i))
		return Cafe{}, err
	}
	value.ID = i
	cs.Lock()
	defer cs.Unlock()
	value = cs.set(i, value)
	return value, nil
}

//=====================HttpResponses======================

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HttpResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []HttpError `json:"errors"`
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

func sendServerError(errorMessage string, w http.ResponseWriter) {
	log.Error().Msgf(errorMessage)
	w.WriteHeader(http.StatusInternalServerError)
}

func sendSingleError(errorMessage string, w http.ResponseWriter) {
	log.Info().Msgf(errorMessage)
	errs := make([]HttpError, 1, 1)
	errs[0] = HttpError{
		Code:    400,
		Message: errorMessage,
	}
	sendSeveralErrors(errs, w)
}

func sendSeveralErrors(errors []HttpError, w http.ResponseWriter) {
	httpResponse := HttpResponse{Errors: errors}
	serializedError, err := json.Marshal(httpResponse)
	if err != nil {
		message := fmt.Sprintf("HttpResponse is json serializing: %s", err.Error())
		sendServerError(message, w)
		return
	}

	_, err = w.Write(serializedError)
	if err != nil {
		message := fmt.Sprintf("HttpResponse while writing is socket: %s", err.Error())
		sendServerError(message, w)
		return
	}
	log.Info().Msgf("Validation error message sent")
}

func sendOKAnswer(data interface{}, w http.ResponseWriter) {
	serializedData, err := json.Marshal(HttpResponse{Data: data})
	if err != nil {
		log.Error().Msgf(err.Error())
		sendServerError("Server JSON encoding error", w)
	}
	_, err = w.Write(serializedData)
	if err != nil {
		message := fmt.Sprintf("HttpResponse while writing is socket: %s", err.Error())
		sendServerError(message, w)
		return
	}
	log.Info().Msgf("OK message sent")
}

//=====================Validator======================

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
		return ut.Add("required", "{{0} is a required field}", true)
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

func getValidationErrors(err error, trans ut.Translator) []HttpError {
	errorsCount := len(err.(validator.ValidationErrors))
	errs := make([]HttpError, errorsCount, errorsCount)

	for i, e := range err.(validator.ValidationErrors) {
		validationError := createNewHttpError(400, e.Translate(trans))
		errs[i] = *validationError
	}
	return errs
}

//=====================Cookies======================

func getAuthCookie(email, password string) (http.Cookie, error) {
	expiresDate := time.Now().Add(time.Hour * 24 * 100)
	token, err := sessions.Login(email, password, expiresDate)

	if err != nil {
		err := errors.New("user with given email and password does not exist")
		return http.Cookie{}, err
	}
	cookie := http.Cookie{
		Name:     "authCookie",
		Value:    token,
		Expires:  expiresDate,
		Path:     "/",
		HttpOnly: true,
	}
	return cookie, nil
}

//=====================Handlers======================

func registerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		sendSingleError("bad request", w)
		return
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		sendSingleError("empty jsonData field", w)
		return
	}

	ownerObj := Owner{}

	if err := json.Unmarshal([]byte(jsonData), &ownerObj); err != nil {
		sendSingleError("json parsing error", w)
		return
	}
	ownerObj.EditedAt = time.Now()

	validation, trans, err := getValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}

	if err := validation.Struct(ownerObj); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := ReceiveFile(file, handler, "owners")
		if err == nil {
			ownerObj.Photo = fmt.Sprintf("%s/%s", serverUrl, filename)
		}
	}

	err, owner := owners.Append(ownerObj)
	if err != nil {
		sendSingleError("User with this email already existed", w)
		return
	}

	cookie, err := getAuthCookie(ownerObj.Email, ownerObj.Password)
	if err != nil {
		message := fmt.Sprintf("troubles with cookies %s", err)
		log.Error().Msgf(message)
		return
	}
	http.SetCookie(w, &cookie)

	sendOKAnswer(owner, w)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		message := fmt.Sprintf("HttpResponse while writing is socket: %s", err.Error())
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
		message := fmt.Sprintf("HttpResponse while unmarshelling: %s", err.Error())
		sendServerError(message, w)
		return
	}

	validation, trans, err := getValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}
	if err := validation.Struct(form); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}

	cookie, err := getAuthCookie(form.Email, form.Password)
	if err != nil {
		sendSingleError("no user with given login and password", w)
		return
	}
	http.SetCookie(w, &cookie)
	sendOKAnswer("", w)
}

func sendForbidden(w http.ResponseWriter) {
	sendSingleError("no permissions", w)
}

func EditOwnerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		sendSingleError("bad request", w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		sendSingleError(message, w)
	}

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

	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		sendSingleError("empty jsonData field", w)
		return
	}
	ownerObj := Owner{}

	if err := json.Unmarshal([]byte(jsonData), &ownerObj); err != nil {
		sendSingleError("json parsing error", w)
		return
	}
	ownerObj.EditedAt = time.Now()

	validation, trans, err := getValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}

	if err := validation.Struct(ownerObj); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := ReceiveFile(file, handler, "owners")
		if err == nil {
			ownerObj.Photo = fmt.Sprintf("%s/%s", serverUrl, filename)
		}
	}

	owner, err = owners.Set(id, ownerObj)
	if err != nil {
		sendSingleError(err.Error(), w)
		return
	}
	sendOKAnswer(owner, w)
}

func getOwnerHandler(w http.ResponseWriter, r *http.Request) {

	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		sendForbidden(w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		sendSingleError(message, w)
	}

	owner, err := owners.Get(id)
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

func createCafeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		sendSingleError("bad request", w)
		return
	}

	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		sendForbidden(w)
		return
	}

	owner, err := sessions.getOwnerByCookie(authCookie.Value)
	if err != nil {
		sendForbidden(w)
		return
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		sendSingleError("empty jsonData field", w)
		return
	}
	cafeObj := Cafe{OwnerID: owner.ID}

	if err := json.Unmarshal([]byte(jsonData), &cafeObj); err != nil {
		sendSingleError("json parsing error", w)
		return
	}

	validation, trans, err := getValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}
	if err := validation.Struct(cafeObj); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := ReceiveFile(file, handler, "cafes")
		if err == nil {
			cafeObj.Photo = fmt.Sprintf("%s/%s", serverUrl, filename)
		}
	}

	err, cafe := cafes.Append(cafeObj)
	if err != nil {
		sendSingleError("user with this email already existed", w)
		return
	}

	sendOKAnswer(cafe, w)
	return
}

func getCurrentOwnerHandler(w http.ResponseWriter, r *http.Request) {
	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		sendForbidden(w)
		return
	}
	owner, err := sessions.getOwnerByCookie(authCookie.Value)
	if err != nil {
		sendForbidden(w)
		return
	}
	sendOKAnswer(owner, w)
}

func getCafesListHandler(w http.ResponseWriter, r *http.Request) {
	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		sendForbidden(w)
		return
	}
	owner, err := sessions.getOwnerByCookie(authCookie.Value)
	if err != nil {
		sendForbidden(w)
		return
	}
	ownerCafes := cafes.getOwnerCafes(owner)
	sendOKAnswer(ownerCafes, w)
}

func getCafeHandler(w http.ResponseWriter, r *http.Request) {
	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		sendForbidden(w)
		return
	}

	owner, err := sessions.getOwnerByCookie(authCookie.Value)
	if err != nil {
		sendForbidden(w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		sendSingleError(message, w)
	}

	cafe, err := cafes.Get(id)
	if err != nil {
		sendForbidden(w)
		return
	}

	if !cafe.hasPermission(owner) {
		sendForbidden(w)
		return
	}
	sendOKAnswer(cafe, w)
}

func EditCafeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		sendSingleError("bad request", w)
		return
	}

	authCookie, err := r.Cookie("authCookie")
	if err != nil {
		sendForbidden(w)
		return
	}

	owner, err := sessions.getOwnerByCookie(authCookie.Value)
	if err != nil {
		sendForbidden(w)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		sendSingleError(message, w)
	}

	cafeObj, err := cafes.Get(id)
	if err != nil {
		sendForbidden(w)
		return
	}

	if !cafeObj.hasPermission(owner) {
		sendForbidden(w)
		return
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" {
		sendSingleError("empty jsonData field", w)
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &cafeObj); err != nil {
		sendSingleError("json parsing error", w)
		return
	}

	validation, trans, err := getValidator()
	if err != nil {
		message := fmt.Sprintf("HttpResponse in validator: %s", err.Error())
		sendServerError(message, w)
		return
	}

	if err := validation.Struct(cafeObj); err != nil {
		errs := getValidationErrors(err, trans)
		sendSeveralErrors(errs, w)
		return
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := ReceiveFile(file, handler, "cafes")
		if err == nil {
			cafeObj.Photo = fmt.Sprintf("%s/%s", serverUrl, filename)
		}
	}

	cafeObj, err = cafes.Set(id, cafeObj)
	if err != nil {
		sendSingleError(err.Error(), w)
		return
	}
	sendOKAnswer(cafeObj, w)
}

func ReceiveFile(file multipart.File, handler *multipart.FileHeader, folder string) (string, error) {

	defer file.Close()

	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	uString := u.String()
	folderName := []rune(uString)[:3]
	separatedFilename := strings.Split(handler.Filename, ".")
	if len(separatedFilename) <= 1 {
		err := errors.New("bad filename")
		return "", err
	}
	fileType := separatedFilename[len(separatedFilename)-1]

	path := fmt.Sprintf("%s/%s/%s", mediaFolder, folder, string(folderName))
	filename := fmt.Sprintf("%s.%s", uString, fileType)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", nil
	}

	fullFilename := fmt.Sprintf("%s/%s", path, filename)

	f, err := os.OpenFile(fullFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return fullFilename, err
}

//=====================Middleware======================

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("URL: %s, METHOD: %s", r.RequestURI, r.Method)
		log.Info().Msgf(msg)
		next.ServeHTTP(w, r)
	})
}

func MyCORSMethodMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "*")
			w.Header().Set("Access-Control-Allow-Methods",
				"POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length,"+
				" Accept-Encoding, X-CSRF-Token, csrf-token, Authorization")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Set-Cookie", "*")
			w.Header().Set("Vary", "Accept, Cookie")
			if req.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

//=====================Storage======================

var owners = NewOwnersStorage()
var sessions = NewSessionsStorage()
var cafes = NewCafesStorage()

//=====================Static settings======================

const mediaFolder = "media"
const serverUrl = "http://localhost:8080"

func main() {
	r := mux.NewRouter()

	//Middleware
	r.Use(MyCORSMethodMiddleware(r))
	r.Use(loggingMiddleware)

	//owner handlers
	r.HandleFunc("/api/v1/owner", registerHandler).Methods("POST")
	r.HandleFunc("/api/v1/owner/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/v1/owner/{id:[0-9]+}", getOwnerHandler).Methods("GET")
	r.HandleFunc("/api/v1/getCurrentOwner/", getCurrentOwnerHandler).Methods("GET")
	r.HandleFunc("/api/v1/owner/{id:[0-9]+}", EditOwnerHandler).Methods("PUT")

	//cafe handlers
	r.HandleFunc("/api/v1/cafe", createCafeHandler).Methods("POST")
	r.HandleFunc("/api/v1/cafe", getCafesListHandler).Methods("GET")
	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}", getCafeHandler).Methods("GET")
	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}", EditCafeHandler).Methods("PUT")

	//OPTIONS
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length,"+
			" Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers,"+
			" Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	//static server
	r.PathPrefix("/media/").Handler(http.StripPrefix("/media/", http.FileServer(http.Dir(mediaFolder))))

	http.Handle("/", r)
	log.Info().Msgf("starting server at :8080")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Error().Msgf(srv.ListenAndServe().Error())

}
