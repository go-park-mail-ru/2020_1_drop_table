package http

import (
	"2020_1_drop_table/internal/app/addStaff"
)

type AddStaffHandler struct {
	AddStaffUsecase addStaff.Usecase
}
