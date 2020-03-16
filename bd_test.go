package main

import (
	"2020_1_drop_table/owners"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAppend(t *testing.T) {
	Storage, _ := owners.NewStaffStorage("postgres", "", "5431")
	Storage.Clear()
	stf := owners.Staff{
		Name:     "asd",
		Email:    "asd",
		Password: "asd",
		EditedAt: time.Now(),
		Photo:    "asd",
	}
	_, err := Storage.Append(stf)
	fmt.Println(err)
	assert.Nil(t, err, "No errors")
	own2 := owners.Staff{
		Name:     "asd",
		Email:    "assd",
		Password: "asd",
		EditedAt: time.Now(),
		Photo:    "asd",
	}
	_, err2 := Storage.Append(own2)

	assert.Nil(t, err2, "No erros")
}

func TestGet(t *testing.T) {
	Storage, _ := owners.NewStaffStorage("postgres", "", "5431")
	Storage.Clear()
	stf := owners.Staff{
		StaffID:  1,
		Name:     "asd",
		Email:    "asd",
		Password: "asd",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	_, _ = Storage.Append(stf)
	dbStaff, err := Storage.Get(1)
	assert.Nil(t, err)
	stf.Password = owners.GetMD5Hash(stf.Password)
	assert.Equal(t, stf, dbStaff)
}

func TestStaffStorage_GetByEmailAndPassword(t *testing.T) {
	Storage, _ := owners.NewStaffStorage("postgres", "", "5431")
	Storage.Clear()
	stf := owners.Staff{
		StaffID:  1,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	_, _ = Storage.Append(stf)
	stf.Password = owners.GetMD5Hash(stf.Password)
	dbStaff, _, err := Storage.GetByEmailAndPassword("email", "password")
	assert.Nil(t, err)
	assert.Equal(t, stf, dbStaff)
}

func TestStaffStorage_Set(t *testing.T) {
	Storage, _ := owners.NewStaffStorage("postgres", "", "5431")
	Storage.Clear()
	stf := owners.Staff{
		StaffID:  1,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}

	newOwn := owners.Staff{
		StaffID:  1,
		Name:     "newasd",
		Email:    "newemail",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	_, _ = Storage.Append(stf)
	_, _ = Storage.Set(1, newOwn)
	dbStaff, err := Storage.Get(1)
	assert.Nil(t, err)
	newOwn.Password = owners.GetMD5Hash(newOwn.Password)
	assert.Equal(t, newOwn, dbStaff)
}

func TestStaffStorage_Count(t *testing.T) {
	Storage, _ := owners.NewStaffStorage("postgres", "", "5431")
	Storage.Clear()
	count, err := Storage.Count()
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	stf := owners.Staff{
		StaffID:  229,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	_, _ = Storage.Append(stf)

	count, err = Storage.Count()
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

}

func TestStaffStorage_Existed(t *testing.T) {
	Storage, _ := owners.NewStaffStorage("postgres", "", "5431")
	Storage.Clear()
	_, isExist, err := Storage.GetByEmailAndPassword("email", "password")
	assert.Nil(t, err)
	assert.Equal(t, false, isExist)

	stf := owners.Staff{
		StaffID:  229,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	_, _ = Storage.Append(stf)

	_, isExist, err = Storage.GetByEmailAndPassword("email", "password")
	assert.Nil(t, err)
	assert.Equal(t, true, isExist)
}
