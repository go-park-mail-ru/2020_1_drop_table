package test

import (
	"2020_1_drop_table/configs"
	rep "2020_1_drop_table/internal/app/staff/repository"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestGetStaffByOwner(t *testing.T) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s",
		configs.PostgresPreferences.User,
		configs.PostgresPreferences.Password,
		configs.PostgresPreferences.Port)

	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	repo := rep.NewPostgresStaffRepository(conn)
	res, err := repo.GetStaffListByOwnerId(context.TODO(), 10)
	fmt.Println(err)
	assert.Equal(t, res, "asd")
}
