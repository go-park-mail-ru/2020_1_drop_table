package http

import (
	"2020_1_drop_table/internal/app/customer"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/pkg/permissions"
	"2020_1_drop_table/internal/pkg/responses"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type CustomerHandler struct {
	CUsecase customer.Usecase
}

func NewCustomerHandler(r *mux.Router, us customer.Usecase) {
	handler := CustomerHandler{CUsecase: us}
	r.HandleFunc("/api/v1/customers/points-system/{uuid}/points/", permissions.SetCSRF(handler.GetPoints)).Methods("GET")
	r.HandleFunc("/api/v1/customers/points-system/{uuid}/{points:[0-9]+}/", permissions.CheckCSRF(handler.SetPoints)).Methods("PUT")
	r.HandleFunc("/api/v1/customers/sale-system/{uuid}", permissions.CheckCSRF(handler.GetSale)).Methods("PUT")

}

func (h CustomerHandler) GetPoints(writer http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	points, err := h.CUsecase.GetPoints(r.Context(), uuid)
	if err != nil {
		responses.SendSingleError(globalModels.ErrBadUuid.Error(), writer)
		return
	}
	responses.SendOKAnswer(points, writer)
}

func (h CustomerHandler) SetPoints(writer http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	newPoints, err := strconv.Atoi(mux.Vars(r)["points"])
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
	}
	err = h.CUsecase.SetPoints(r.Context(), uuid, newPoints)
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
		return
	}
	responses.SendOKAnswer(newPoints, writer)
}

func FetchSum(r *http.Request) (float32, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return 0, globalModels.ErrBadRequest
	}

	jsonData := r.FormValue("jsonData")
	if jsonData == "" || jsonData == "null" {
		return 0, globalModels.ErrEmptyJSON
	}

	var sum float32

	if err := json.Unmarshal([]byte(jsonData), &sum); err != nil {
		return 0, globalModels.ErrBadJSON
	}
	return sum, err
}

func (h CustomerHandler) GetSale(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	sum, err := FetchSum(r)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
	}
	newSum, err := h.CUsecase.GetSale(r.Context(), sum, uuid)
	if err != nil {
		responses.SendSingleError(err.Error(), w)
		return
	}
	responses.SendOKAnswer(newSum, w)
}
