package http

import (
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/statistics"
	models2 "2020_1_drop_table/internal/app/statistics/models"
	"2020_1_drop_table/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type StatisticsHandler struct {
	SUsecase statistics.Usecase
}

func NewStatisticsHandler(r *mux.Router, us statistics.Usecase) {
	handler := StatisticsHandler{
		SUsecase: us,
	}

	r.HandleFunc("/api/v1/statistics/get_worker_data", handler.GetWorkerData).Methods("POST") //todo check csrf
	r.HandleFunc("/api/v1/statistics/get_graphs_data", handler.GetGraphsData).Methods("GET")  //todo check csrf

}

func fetchWorkerData(r *http.Request) (models2.GetWorkerDataStruct, error) {

	data, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		return models2.GetWorkerDataStruct{}, err
	}

	if len(data) == 0 {
		return models2.GetWorkerDataStruct{}, globalModels.ErrBadJSON
	}

	var WorkerData models2.GetWorkerDataStruct
	err = json.Unmarshal(data, &WorkerData)
	fmt.Println(err)
	if err != nil {
		return models2.GetWorkerDataStruct{}, globalModels.ErrBadJSON
	}
	return WorkerData, nil
}

func (h StatisticsHandler) GetWorkerData(writer http.ResponseWriter, request *http.Request) {
	//todo permissions only for this cafe staff
	workerData, err := fetchWorkerData(request)
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
		return
	}
	res, err := h.SUsecase.GetWorkerData(request.Context(), workerData.StaffID, workerData.Limit, workerData.Since)
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
		return
	}
	responses.SendOKAnswer(res, writer)
}

func (h StatisticsHandler) GetGraphsData(writer http.ResponseWriter, request *http.Request) {
	typ := request.FormValue("type")
	since := request.FormValue("since")
	to := request.FormValue("to")
	fmt.Println(typ, since, to)
	jsonData, err := h.SUsecase.GetDataForGraphs(request.Context(), typ, since, to)
	if err != nil {
		responses.SendSingleError(err.Error(), writer)
		return
	}
	responses.SendOKAnswer(jsonData, writer)

}
