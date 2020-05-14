package usecase

import (
	cafeClient "2020_1_drop_table/internal/app/cafe/delivery/grpc/client"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/statistics"
	"2020_1_drop_table/internal/app/statistics/models"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type statisticsUsecase struct {
	statisticsRepo statistics.Repository
	staffClient    staffClient.StaffClientInterface
	cafeClient     cafeClient.CafeGRPCClientInterface
	contextTimeout time.Duration
}

func (s statisticsUsecase) AddData(jsonData string, time time.Time, clientUUID string, staffId int, cafeId int) error {
	return s.statisticsRepo.AddData(jsonData, time, clientUUID, staffId, cafeId)

}

func (s statisticsUsecase) GetWorkerData(ctx context.Context, staffID int, limit int, since int) ([]models.StatisticsStruct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	requestStaff, err := s.staffClient.GetFromSession(ctx)

	//todo permissions check if staffId working on request user
	if err != nil || !requestStaff.IsOwner {
		message := "Request user not owner or not authorized"
		return nil, errors.New(message)
	}

	data, err := s.statisticsRepo.GetWorkerDataFromRepo(ctx, staffID, limit, since)

	return data, err
}

func NewStatisticsUsecase(st statistics.Repository, staffClient staffClient.StaffClientInterface, cafeClient cafeClient.CafeGRPCClientInterface, timeout time.Duration) *statisticsUsecase {
	return &statisticsUsecase{
		statisticsRepo: st,
		contextTimeout: timeout,
		staffClient:    staffClient,
		cafeClient:     cafeClient,
	}
}

func (s statisticsUsecase) GetDataForGraphs(ctx context.Context, typ string, since string, to string) (map[string]map[string][]models.TempStruct, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	requestUser, err := s.staffClient.GetFromSession(ctx)
	if err != nil {
		return nil, globalModels.ErrForbidden
	}
	cafes, err := s.cafeClient.GetByOwnerId(ctx, requestUser.StaffID)
	if err != nil {
		message := fmt.Sprintf("cafe not found")
		return nil, errors.New(message)
	}
	rawData, err := s.statisticsRepo.GetGraphsDataFromRepo(ctx, cafes, typ, since, to)
	if err != nil {
		message := fmt.Sprintf("cant found statistics data with this input")
		return nil, errors.New(message)
	}
	fmt.Println(rawData)
	jsonData := jsonify(rawData)
	return jsonData, nil
}

func jsonify(data []models.StatisticsGraphRawStruct) map[string]map[string][]models.TempStruct {
	//todo refactor work super slow

	var m = make(map[string][]models.TempStruct)
	var resMap = make(map[string]map[string][]models.TempStruct)

	for _, el := range data {
		staffIdAndCafe := strconv.Itoa(el.StaffId) + "_" + strconv.Itoa(el.CafeId)
		m[staffIdAndCafe] = append(m[staffIdAndCafe], models.TempStruct{
			NumOfUsage: el.Count,
			Date:       el.Date,
		})
	}
	tempMap := make(map[string][]models.TempStruct)
	prevKey := getFirstKeyFromMap(m)
	fmt.Println(prevKey)
	for key, value := range m {
		keysArr := strings.Split(key, "_")
		staffId := keysArr[0]
		cafeId := keysArr[1]
		if prevKey != cafeId {
			tempMap = make(map[string][]models.TempStruct)
		}
		tempMap[staffId] = append(tempMap[staffId], value...)
		resMap[cafeId] = tempMap
		prevKey = cafeId
	}
	return resMap
}

func getFirstKeyFromMap(m map[string][]models.TempStruct) string {
	for key := range m {
		return strings.Split(key, `_`)[1]
	}
	return ""
}
