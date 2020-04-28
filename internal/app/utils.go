package app

import (
	"2020_1_drop_table/configs"
	staffModels "2020_1_drop_table/internal/microservices/staff/models"
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func GetSafeStaff(unsafeStaff staffModels.Staff) staffModels.SafeStaff {
	return staffModels.SafeStaff{
		StaffID:  unsafeStaff.StaffID,
		Name:     unsafeStaff.Name,
		Email:    unsafeStaff.Email,
		EditedAt: unsafeStaff.EditedAt,
		Photo:    unsafeStaff.Photo,
		IsOwner:  unsafeStaff.IsOwner,
		CafeId:   unsafeStaff.CafeId,
		Position: unsafeStaff.Position,
	}
}

func SaveFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {

	defer file.Close()

	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	uString := u.String()
	folderName := []rune(uString)[:3]
	separatedFilename := strings.Split(header.Filename, ".")
	if len(separatedFilename) <= 1 {
		err := errors.New("bad filename")
		return "", err
	}
	fileType := separatedFilename[len(separatedFilename)-1]

	path := fmt.Sprintf("%s/%s/%s", configs.MediaFolder, folder, string(folderName))
	filename := fmt.Sprintf("%s.%s", uString, fileType)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", nil
	}

	fullFilename := fmt.Sprintf("%s/%s", path, filename)

	f, err := os.OpenFile(fullFilename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	return fullFilename, err
}
