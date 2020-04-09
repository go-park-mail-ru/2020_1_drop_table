package models

type Customer struct {
	CustomerID string `structs:"CustomerID"`
	CafeID     int    `structs:"-"`
	Points     int    `structs:"Points"`
}
