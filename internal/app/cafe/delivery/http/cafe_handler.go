package http

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app"
	"2020_1_drop_table/internal/app/cafe"
	"2020_1_drop_table/internal/app/cafe/models"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/pkg/permissions"
	"2020_1_drop_table/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type cafeHandler struct {
	CUsecase cafe.Usecase
}

func NewCafeHandler(r *mux.Router, us cafe.Usecase) {
	handler := cafeHandler{
		CUsecase: us,
	}

	r.HandleFunc("/api/v1/cafe", permissions.CheckCSRF(permissions.CheckAuthenticated(handler.AddCafeHandler))).Methods("POST")
	r.HandleFunc("/api/v1/cafe", permissions.SetCSRF(handler.GetByOwnerIDHandler)).Methods("GET")
	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}", permissions.SetCSRF(handler.GetByIDHandler)).Methods("GET")
	r.HandleFunc("/api/v1/cafe/{id:[0-9]+}", permissions.CheckCSRF(permissions.CheckAuthenticated(handler.EditCafeHandler))).Methods("PUT")
}

func (c *cafeHandler) fetchCafe(r *http.Request) (models.Cafe, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return models.Cafe{}, globalModels.ErrBadRequest
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" || jsonData == "null" {
		return models.Cafe{}, globalModels.ErrEmptyJSON
	}

	cafeObj := models.Cafe{}

	if err := json.Unmarshal([]byte(jsonData), &cafeObj); err != nil {
		return models.Cafe{}, globalModels.ErrBadJSON
	}
	if file, handler, err := r.FormFile("photo"); err == nil {
		filename, err := app.SaveFile(file, handler, "cafe")
		if err == nil {
			cafeObj.Photo = fmt.Sprintf("%s/%s", configs.ServerUrl, filename)
		}
	}

	return cafeObj, nil
}

func (c *cafeHandler) AddCafeHandler(w http.ResponseWriter, r *http.Request) {
	cafeObj, err := c.fetchCafe(r)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	cafeObj, err = c.CUsecase.Add(r.Context(), cafeObj)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(cafeObj, w)
	return
}

func (c *cafeHandler) EditCafeHandler(w http.ResponseWriter, r *http.Request) {
	cafeObj, err := c.fetchCafe(r)
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

	cafeObj.CafeID = id

	err = c.CUsecase.Update(r.Context(), cafeObj)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(cafeObj, w)
	return
}

func (c *cafeHandler) GetByOwnerIDHandler(w http.ResponseWriter, r *http.Request) {
	cafesObj, err := c.CUsecase.GetByOwnerID(r.Context())
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(cafesObj, w)
	return
}

func (c *cafeHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		message := fmt.Sprintf("bad id: %s", mux.Vars(r)["id"])
		responses.SendSingleError(message, w)
		return
	}

	cafeObj, err := c.CUsecase.GetByID(r.Context(), id)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}

	responses.SendOKAnswer(cafeObj, w)
	return
}
