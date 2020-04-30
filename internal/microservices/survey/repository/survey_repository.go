package repository

import (
	"2020_1_drop_table/internal/microservices/survey"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgresSurveyRepository struct {
	Conn *sqlx.DB
}

func NewPostgresSurveyRepository(conn *sqlx.DB) survey.Repository {
	return &postgresSurveyRepository{
		Conn: conn,
	}
}

func (p postgresSurveyRepository) SetSurveyTemplate(ctx context.Context, survey string, id int, cafeOwnerID int) error {
	query := `INSERT INTO surveytemplate (cafeid, surveytemplate,cafeOwnerId) VALUES ($1,$2,$3)`
	_, err := p.Conn.ExecContext(ctx, query, id, survey, cafeOwnerID)
	return err
}
func (p postgresSurveyRepository) GetSurveyTemplate(ctx context.Context, cafeId int) (string, error) {
	query := `SELECT surveytemplate from surveytemplate where cafeid=$1`
	var survTemplate string
	err := p.Conn.GetContext(ctx, &survTemplate, query, cafeId)
	return survTemplate, err
}

func (p postgresSurveyRepository) SubmitSurvey(ctx context.Context, surv string, customerUUID string) error {
	query := `UPDATE customer SET surveyresult=$1 where customerID=$2`
	_, err := p.Conn.ExecContext(ctx, query, surv, customerUUID)
	return err
}
func (p postgresSurveyRepository) UpdateSurveyTemplate(ctx context.Context, survey string, id int) error {
	query := `UPDATE surveytemplate set surveytemplate=$1 where cafeid=$2`
	_, err := p.Conn.ExecContext(ctx, query, survey, id)
	return err
}
