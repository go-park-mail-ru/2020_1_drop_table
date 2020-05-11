package models

type GetWorkerDataStruct struct {
	StaffID int `json:"staffID"`
	Since   int `json:"since"`
	Limit   int `json:"limit"`
}
