package survey

import "context"

type Usecase interface {
	SetSurveyTemplate(ctx context.Context, survey string, id int) error
	GetSurveyTemplate(ctx context.Context, id int) (string, error)
	SubmitSurvey(ctx context.Context, survey string, customerID string) error
}
