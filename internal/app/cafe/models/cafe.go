package models

import (
	"database/sql"
	"time"
)

type Cafe struct {
	CafeID               int           `json:"id"`
	Name                 string        `json:"name" validate:"required,min=2,max=100"`
	Address              string        `json:"address" validate:"required"`
	Description          string        `json:"description" validate:"required"`
	StaffID              int           `json:"staffID"`
	OpenTime             time.Time     `json:"openTime"`
	CloseTime            time.Time     `json:"closeTime"`
	Photo                string        `json:"photo"`
	PublishedApplePassID sql.NullInt64 `json:"-"`
	SavedApplePassID     sql.NullInt64 `json:"-"`
}
