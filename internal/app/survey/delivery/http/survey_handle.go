package http

import (
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/survey"
	"2020_1_drop_table/internal/pkg/permissions"
	"2020_1_drop_table/internal/pkg/responses"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type SurveyHandler struct {
	SurveyUC survey.Usecase
}

func NewSurveyHandler(r *mux.Router, us survey.Usecase) {
	handler := SurveyHandler{
		SurveyUC: us,
	}
	r.HandleFunc("/api/v1/survey/set_survey_template/{id:[0-9]+}", permissions.CheckCSRF(handler.SetSurveyTemplate)).Methods("POST")
	r.HandleFunc("/api/v1/survey/get_survey_template/{id:[0-9]+}", permissions.SetCSRF(handler.GetSurveyTemplate)).Methods("GET")
	r.HandleFunc("/api/v1/survey/submit_survey/{customerid}", permissions.CheckCSRF(handler.SubmitSurvey)).Methods("POST")

}

func fetchSurvey(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", globalModels.ErrBadRequest
	}
	jsonData := r.FormValue("jsonData")
	return jsonData, nil
}

func (c *SurveyHandler) SetSurveyTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}
	surv, err := fetchSurvey(r)
	if err != nil {
		message := fmt.Sprintf("bad json data")
		responses.SendSingleError(message, w)
		return
	}
	err = c.SurveyUC.SetSurveyTemplate(r.Context(), surv, id)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(err, w)
}

func (c *SurveyHandler) GetSurveyTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}
	survTemplate, err := c.SurveyUC.GetSurveyTemplate(r.Context(), id)
	if err != nil {
		responses.SendSingleErrorWithMessage(err, globalModels.ErrForbidden.Error(), w)
		return
	}
	responses.SendOKAnswer(survTemplate, w)
}

func (c *SurveyHandler) SubmitSurvey(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["customerid"]
	surv, err := fetchSurvey(r)
	if err != nil {
		message := fmt.Sprintf("bad json data")
		responses.SendSingleError(message, w)
		return
	}
	err = c.SurveyUC.SubmitSurvey(r.Context(), surv, customerID)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(nil, w)

}
