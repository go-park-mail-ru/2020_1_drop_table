package usecase

import (
	"2020_1_drop_table/internal/app/addStaff/repository"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddStaffUsecase_GetQrForStaff(t *testing.T) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s", "postgres", "", "5431")
	db, err := sqlx.Open("postgres", connStr)
	assert.Nil(t, err)
	rep := repository.NewUuidCafeRepository(db)
	usecase := newAddStaffUsecase(&rep, time.Second*10)
	code, err := usecase.GetQrForStaff(context.TODO(), 2)
	assert.Nil(t, err)
	fmt.Println(code)
}
