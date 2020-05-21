package models

type StaffByOwnerResponse struct {
	CafeId   string  `json:"CafeId"`
	CafeName string  `json:"CafeName"`
	StaffId  *int    `json:"StaffId"`
	Photo    *string `json:"Photo"`
	Name     *string `json:"StaffName"`
	Position *string `json:"Position"`
}
