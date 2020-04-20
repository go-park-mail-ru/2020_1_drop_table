package models

type Customer struct {
	CustomerID   string  `structs:"CustomerID"`
	CafeID       int     `structs:"-"`
	Points       int     `structs:"Points"`
	Sum          float32 `structs:"Sum"`
	SurveyResult string  `structs:"SurveyResult"`
}
