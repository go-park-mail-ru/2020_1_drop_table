package models

import "time"

type GetWorkerDataStruct struct {
	StaffID int `json:"staffID"`
	Since   int `json:"since"`
	Limit   int `json:"limit"`
}

type StatisticsStruct struct {
	JsonData   string    `json:"json_data"`
	Time       time.Time `json:"time"`
	ClientUUID string    `json:"clientUuid"`
	StaffId    int       `json:"staffId"`
	CafeId     int       `json:"cafeId"`
}

type StatisticsGraphRawStruct struct {
	Count   int       `json:"count"`
	Date    time.Time `json:"date"`
	CafeId  int       `json:"cafeId"`
	StaffId int       `json:"staffId"`
}

type TempStruct struct {
	NumOfUsage int
	Date       time.Time
}
