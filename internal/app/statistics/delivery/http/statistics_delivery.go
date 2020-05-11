package http

import (
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/statistics"
	models2 "2020_1_drop_table/internal/app/statistics/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type StatisticsHandler struct {
	SUsecase statistics.Usecase
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
	fmt.Println(workerData, err)
}

func NewStatisticsHandler(r *mux.Router, us statistics.Usecase) {
	handler := StatisticsHandler{
		SUsecase: us,
	}

	r.HandleFunc("/api/v1/statistics/get_worker_data", handler.GetWorkerData).Methods("GET")

}
