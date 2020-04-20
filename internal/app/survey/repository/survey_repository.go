package repository

import (
	"2020_1_drop_table/internal/app/survey"
	"context"
	"github.com/jmoiron/sqlx"
)

type postgresSurveyRepository struct {
	Conn *sqlx.DB
}

func NewPostgresCafeRepository(conn *sqlx.DB) survey.Repository {
	return &postgresSurveyRepository{
		Conn: conn,
	}
}

func (p postgresSurveyRepository) SetSurveyTemplate(ctx context.Context, survey string, id int, cafeOwnerID int) error {
	query := `INSERT INTO surveytemplate (cafeid, surveytemplate,cafeOwnerId) VALUES ($1,$2,$3)`
	_, err := p.Conn.ExecContext(ctx, query, id, survey, cafeOwnerID)
	return err
}
func (p postgresSurveyRepository) GetSurveyTemplate(ctx context.Context, cafeId int, cafeOwnerId int) (string, error) {
	query := `SELECT surveytemplate from surveytemplate where cafeid=$1 and cafeownerid=$2`
	var survTemplate string
	err := p.Conn.GetContext(ctx, &survTemplate, query, cafeId, cafeOwnerId)
	return survTemplate, err
}

func (p postgresSurveyRepository) SubmitSurvey(ctx context.Context, surv string, customerUUID string) error {
	query := `UPDATE customer SET surveyresult=$1 where customerID=$2`
	_, err := p.Conn.ExecContext(ctx, query, surv, customerUUID)
	return err
}