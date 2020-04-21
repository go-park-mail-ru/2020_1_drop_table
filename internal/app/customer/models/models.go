package models

type Customer struct {
	CustomerID   string `json:"customer_id" structs:"CustomerID"`
	CafeID       int    `json:"cafe_id" structs:"-"`
	Type         string `json:"type" structs:"-"`
	Points       string `json:"points" structs:"Points"`
	SurveyResult string `json:"survey_result" structs:"SurveyResult"`
}
