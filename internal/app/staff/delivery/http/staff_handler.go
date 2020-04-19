package http

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/staff"
	"2020_1_drop_table/internal/app/staff/models"
	"2020_1_drop_table/internal/pkg/permissions"
	"2020_1_drop_table/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
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
	r.HandleFunc("/api/v1/staff", permissions.SetCSRF(handler.RegisterHandler)).Methods("POST")
	r.HandleFunc("/api/v1/get_current_staff/", permissions.SetCSRF(handler.GetCurrentStaffHandler)).Methods("GET")
	r.HandleFunc("/api/v1/staff/login", permissions.SetCSRF(handler.LoginHandler)).Methods("POST")
	r.HandleFunc("/api/v1/staff/{id:[0-9]+}", permissions.SetCSRF(permissions.CheckAuthenticated(handler.GetStaffByIdHandler))).Methods("GET")
	r.HandleFunc("/api/v1/staff/{id:[0-9]+}", permissions.CheckCSRF(permissions.CheckAuthenticated(handler.EditStaffHandler))).Methods("PUT")
	r.HandleFunc("/api/v1/staff/generateQr/{id:[0-9]+}", permissions.SetCSRF(handler.GenerateQrHandler)).Methods("GET")
	r.HandleFunc("/api/v1/add_staff", permissions.SetCSRF(handler.AddStaffHandler)).Methods("POST")
	r.HandleFunc("/api/v1/staff/get_staff_list/{id:[0-9]+}", permissions.SetCSRF(handler.GetStaffListHandler)).Methods("GET")
	r.HandleFunc("/api/v1/staff/delete_staff/{id:[0-9]+}", permissions.CheckCSRF(handler.DeleteStaff)).Methods("POST")
	r.HandleFunc("/api/v1/staff/update_position/{id:[0-9]+}", permissions.CheckCSRF(handler.UpdatePosition)).Methods("POST") //TODO CHECK CSRF

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
		filename, err := app.SaveFile(file, handler, "staffs")
		if err == nil {
			staffObj.Photo = fmt.Sprintf("%s/%s", configs.ServerUrl, filename)
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

func (s *StaffHandler) AddStaffHandler(w http.ResponseWriter, r *http.Request) {
	staffObj, err := s.fetchStaff(r)
	uuid := r.FormValue("uuid")
	if err != nil && uuid != "" {
		responses.SendSingleError(err.Error(), w)
		return
	}
	CafeId, err := s.SUsecase.GetCafeId(r.Context(), uuid)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	staffObj.IsOwner = false
	staffObj.CafeId = CafeId
	err = s.SUsecase.DeleteQrCodes(uuid)
	if err != nil {
		log.Error().Msgf("error when trying to delete QRCodes")
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

	staffObj, err = s.SUsecase.Update(r.Context(), staffObj)

	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(staffObj, w)
	return
}

func (s *StaffHandler) GenerateQrHandler(w http.ResponseWriter, r *http.Request) {
	CafeId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}
	pathToQr, err := s.SUsecase.GetQrForStaff(r.Context(), CafeId)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(pathToQr, w)
}

func (s *StaffHandler) GetStaffListHandler(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	res, err := s.SUsecase.GetStaffListByOwnerId(r.Context(), ownerId)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(res, w)
}

func (s *StaffHandler) DeleteStaff(w http.ResponseWriter, r *http.Request) {
	staffID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responses.SendSingleError(globalModels.ErrBadRequest.Error(), w)
		return
	}
	err = s.SUsecase.DeleteStaffById(r.Context(), staffID)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(nil, w)
}

func fetchPosition(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", globalModels.ErrBadRequest
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" || jsonData == "null" {
		return "", globalModels.ErrEmptyJSON
	}

	return jsonData, nil
}
func (s *StaffHandler) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	staffID, err := strconv.Atoi(mux.Vars(r)["id"])
	newPosition, fetchErr := fetchPosition(r)
	if err != nil || fetchErr != nil {
		responses.SendSingleError(globalModels.ErrBadRequest.Error(), w)
		return
	}
	err = s.SUsecase.UpdatePosition(r.Context(), staffID, newPosition)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(newPosition, w)

}
