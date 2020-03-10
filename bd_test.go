package main

import (
	"2020_1_drop_table/owners"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAppend(t *testing.T) {
	Storage, _ := newOwnerStorage("postgres", "")
	Storage.PrepareForTest()
	own := Owner{
		OwnerId:  229,
		Name:     "asd",
		Email:    "asd",
		Password: "asd",
		EditedAt: time.Now(),
		Photo:    "asd",
	}
	err := Storage.Append(own)
	assert.Nil(t, err, "No errors")
	own2 := Owner{
		OwnerId:  228,
		Name:     "asd",
		Email:    "asd",
		Password: "asd",
		EditedAt: time.Now(),
		Photo:    "asd",
	}
	err2 := Storage.Append(own2)
	assert.Nil(t, err2, "No erros")
	cantAppend := Storage.Append(own)
	assert.Equal(t, cantAppend.Error(), "pq: duplicate key value violates unique constraint \"owner_pkey\"")
}

func TestGet(t *testing.T) {
	Storage, _ := newOwnerStorage("postgres", "")
	Storage.PrepareForTest()
	own := Owner{
		OwnerId:  229,
		Name:     "asd",
		Email:    "asd",
		Password: "asd",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	Storage.Append(own)
	dbOwner, err := Storage.Get(229)
	assert.Nil(t, err)
	assert.Equal(t, own, dbOwner)
}

func TestOwnerStorage_GetByEmailAndPassword(t *testing.T) {
	Storage, _ := newOwnerStorage("postgres", "")
	Storage.PrepareForTest()
	own := Owner{
		OwnerId:  229,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	Storage.Append(own)

	dbOwner, err := Storage.GetByEmailAndPassword("email", "password")
	assert.Nil(t, err)
	assert.Equal(t, own, dbOwner)
}

func TestOwnerStorage_Set(t *testing.T) {
	Storage, _ := newOwnerStorage("postgres", "")
	Storage.PrepareForTest()
	own := Owner{
		OwnerId:  229,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}

	newOwn := Owner{
		OwnerId:  229,
		Name:     "newasd",
		Email:    "newemail",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	Storage.Append(own)
	Storage.Set(229, newOwn)
	dBOwner, err := Storage.Get(229)
	assert.Nil(t, err)
	assert.Equal(t, newOwn, dBOwner)
}

func TestOwnerStorage_Count(t *testing.T) {
	Storage, _ := newOwnerStorage("postgres", "")
	Storage.PrepareForTest()
	count, err := Storage.Count()
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	own := Owner{
		OwnerId:  229,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	Storage.Append(own)

	count, err = Storage.Count()
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

}

func TestOwnerStorage_AppendList(t *testing.T) {
	Storage, _ := newOwnerStorage("postgres", "")
	Storage.PrepareForTest()

	own := Owner{
		OwnerId:  229,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}

	own2 := Owner{
		OwnerId:  230,
		Name:     "asd2",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}

	ownerList := []Owner{own, own2}
	err := Storage.AppendList(ownerList)

	assert.Nil(t, err)

}

func TestOwnerStorage_Existed(t *testing.T) {
	Storage, _ := newOwnerStorage("postgres", "")
	Storage.PrepareForTest()
	isExist, err := Storage.Existed("email", "password")
	assert.Nil(t, err)
	assert.Equal(t, false, isExist)

	own := Owner{
		OwnerId:  229,
		Name:     "asd",
		Email:    "email",
		Password: "password",
		EditedAt: time.Now().UTC(),
		Photo:    "asd",
	}
	Storage.Append(own)

	isExist, err = Storage.Existed("email", "password")
	assert.Nil(t, err)
	assert.Equal(t, true, isExist)
}
