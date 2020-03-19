package http

import (
	"2020_1_drop_table/app"
	globalModels "2020_1_drop_table/app/models"
	"2020_1_drop_table/app/staff"
	"2020_1_drop_table/app/staff/models"
	"2020_1_drop_table/permissions"
	"2020_1_drop_table/projectConfig"
	"2020_1_drop_table/responses"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"strconv"
)

type StaffHandler struct {
	SUsecase staff.Usecase
}

func NewStaffHandler(r *mux.Router, us staff.Usecase) {
	handler := StaffHandler{
		SUsecase: us,
	}

	r.HandleFunc("/api/v1/staff", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/get_current_staff/", handler.GetCurrentStaffHandler).Methods("GET")
	r.HandleFunc("/api/v1/staff/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/api/v1/staff/{id:[0-9]+}", permissions.CheckAuthenticated(handler.GetStaffByIdHandler)).Methods("GET")
	r.HandleFunc("/api/v1/staff/{id:[0-9]+}", permissions.CheckAuthenticated(handler.EditStaffHandler)).Methods("PUT")

}
func (s *StaffHandler) fetchStaff(r *http.Request) (models.Staff, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return models.Staff{}, globalModels.ErrBadRequest
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" || jsonData == "null" {
		return models.Staff{}, globalModels.ErrEmptyJSON
	}

	var staffObj models.Staff

	if err := json.Unmarshal([]byte(jsonData), &staffObj); err != nil {
		return models.Staff{}, globalModels.ErrBadJSON
	}

	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := s.SUsecase.SaveFile(file, handler, "staffs")
		if err == nil {
			staffObj.Photo = fmt.Sprintf("%s/%s", projectConfig.ServerUrl, filename)
		}
	}

	return staffObj, nil
}

func (s *StaffHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	staffObj, err := s.fetchStaff(r)

	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	safeStaff, err := s.SUsecase.Add(r.Context(), staffObj)

	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)

	session.Values["userID"] = safeStaff.StaffID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.SendOKAnswer(safeStaff, w)
}

func (s *StaffHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	var form models.LoginForm
	err = json.Unmarshal(data, &form)
	if err != nil {
		message := fmt.Sprintf("HttpResponse while serializing: %s", err.Error())
		responses.SendServerError(message, w)
		return
	}

	safeStaff, err := s.SUsecase.GetByEmailAndPassword(r.Context(), form)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	session := r.Context().Value("session").(*sessions.Session)

	session.Values["userID"] = safeStaff.StaffID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.SendOKAnswer(safeStaff, w)
	return
}

func (s *StaffHandler) GetStaffByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	safeStaff, err := s.SUsecase.GetByID(r.Context(), id)

	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(safeStaff, w)
	return
}

func (s *StaffHandler) GetCurrentStaffHandler(w http.ResponseWriter, r *http.Request) {
	staffObj, err := s.SUsecase.GetFromSession(r.Context())

	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(staffObj, w)
	return
}

func (s *StaffHandler) EditStaffHandler(w http.ResponseWriter, r *http.Request) {
	staffUnsafe, err := s.fetchStaff(r)
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

	staffObj := app.GetSafeStaff(staffUnsafe)
	staffObj.StaffID = id

	err = s.SUsecase.Update(r.Context(), staffObj)

	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(staffObj, w)
	return
}
