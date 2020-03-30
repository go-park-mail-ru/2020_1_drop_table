package http

import (
	"2020_1_drop_table/internal/app/customer"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/pkg/responses"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomerHandler struct {
	CUsecase customer.Usecase
}

func NewCustomerHandler(r *mux.Router, us customer.Usecase) {
	handler := CustomerHandler{CUsecase: us}
	r.HandleFunc("/api/v1/customers/{uuid}/points/", handler.GetPoints).Methods("GET")
}

func (h CustomerHandler) GetPoints(writer http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	points, err := h.CUsecase.GetPoints(r.Context(), uuid)
	if err != nil {
		responses.SendSingleError(globalModels.BadUuid.Error(), writer)
		return
	}
	responses.SendOKAnswer(points, writer)

}
