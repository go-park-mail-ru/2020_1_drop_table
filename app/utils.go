package app

import "2020_1_drop_table/app/staff/models"

func GetSafeStaff(unsafeStaff models.Staff) models.SafeStaff {
	return models.SafeStaff{
		StaffID:  unsafeStaff.StaffID,
		Name:     unsafeStaff.Name,
		Email:    unsafeStaff.Email,
		EditedAt: unsafeStaff.EditedAt,
		Photo:    unsafeStaff.Photo,
		IsOwner:  unsafeStaff.IsOwner,
	}
}
