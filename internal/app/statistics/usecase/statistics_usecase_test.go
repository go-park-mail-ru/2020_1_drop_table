package usecase

import (
	"2020_1_drop_table/configs"
	"2020_1_drop_table/internal/app/statistics/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s",
		configs.PostgresPreferences.User,
		configs.PostgresPreferences.Password,
		configs.PostgresPreferences.Port)

	conn, err := sqlx.Open("postgres", connStr)
	fmt.Println(err)
	rep := repository.NewPostgresStatisticsRepository(conn)
	timeoutContext := configs.Timeouts.ContextTimeout

	stUcase := NewStatisticsUsecase(rep, nil, timeoutContext)
	err = stUcase.AddData(`{"asd":"bsd"}`, time.Now(), "c81ccfda-68fc-4cef-8e7a-abc92958e6ee", 4, 2)
	assert.Nil(t, err)

}
