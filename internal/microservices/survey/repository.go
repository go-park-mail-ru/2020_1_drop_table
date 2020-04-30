package survey

import "context"

type Repository interface {
	SetSurveyTemplate(ctx context.Context, survey string, id int, cafeOwnerID int) error
	GetSurveyTemplate(ctx context.Context, cafeId int) (string, error)
	SubmitSurvey(ctx context.Context, surv string, customerID string) error
	UpdateSurveyTemplate(ctx context.Context, survey string, id int) error
}
