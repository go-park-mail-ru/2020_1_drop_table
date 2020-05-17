package models

import (
	"time"
)

//easyjson:json
type Cafe struct {
	CafeID      int       `json:"id"`
	CafeName    string    `json:"name" validate:"required,min=2,max=100"`
	Address     string    `json:"address" validate:"required"`
	Description string    `json:"description" validate:"required"`
	StaffID     int       `json:"staffID"`
	OpenTime    time.Time `json:"openTime"`
	CloseTime   time.Time `json:"closeTime"`
	Photo       string    `json:"photo"`
	Location    string    `json:"location" db:"location_str"`
}

type CafeWithGeoData struct {
	CafeID      int       `json:"id"`
	CafeName    string    `json:"name" validate:"required,min=2,max=100"`
	Address     string    `json:"address" validate:"required"`
	Description string    `json:"description" validate:"required"`
	StaffID     int       `json:"staffID"`
	OpenTime    time.Time `json:"openTime"`
	CloseTime   time.Time `json:"closeTime"`
	Photo       string    `json:"photo"`
	Location    string    `json:"location" db:"location_str"`
}
