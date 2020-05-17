package http

import (
	"2020_1_drop_table/internal/app/customer"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/pkg/permissions"
	"2020_1_drop_table/internal/pkg/responses"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type CustomerHandler struct {
	CUsecase customer.Usecase
}

func NewCustomerHandler(r *mux.Router, us customer.Usecase) {
	handler := CustomerHandler{CUsecase: us}
	r.HandleFunc("/api/v1/customers/{uuid}/points/", permissions.SetCSRF(handler.GetPoints)).Methods("GET")
	r.HandleFunc("/api/v1/customers/{uuid}/customer/", permissions.SetCSRF(handler.GetCustomer)).Methods("GET")
	r.HandleFunc("/api/v1/customers/{uuid}/", handler.SetPoints).Methods("PUT")

}

func (h CustomerHandler) GetCustomer(writer http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	points, err := h.CUsecase.GetCustomer(r.Context(), uuid)
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
		return
	}
	responses.SendOKAnswer(points, writer)
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
		return
	}
	err = h.CUsecase.SetPoints(r.Context(), uuid, string(body))
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
		return
	}
	responses.SendOKAnswer(string(body), writer)
}
