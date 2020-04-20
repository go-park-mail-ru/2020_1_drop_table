package usecase

import (
	"2020_1_drop_table/internal/app/cafe"
	globalModels "2020_1_drop_table/internal/app/models"
	"2020_1_drop_table/internal/app/staff"
	"2020_1_drop_table/internal/app/survey"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SurveyUsecase struct {
	cafeRepo       cafe.Repository
	surveyRepo     survey.Repository
	staffUcase     staff.Usecase
	contextTimeout time.Duration
}

func NewSurveyUsecase(c cafe.Repository, surveyRepo survey.Repository, s staff.Usecase, timeout time.Duration) survey.Usecase {
	return &SurveyUsecase{
		cafeRepo:       c,
		staffUcase:     s,
		surveyRepo:     surveyRepo,
		contextTimeout: timeout,
	}
}

func (s SurveyUsecase) SetSurveyTemplate(ctx context.Context, survey string, id int) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	requestUser, err := s.staffUcase.GetFromSession(ctx)
	if err != nil || !requestUser.IsOwner {
		return globalModels.ErrForbidden
	}
	caf, err := s.cafeRepo.GetByID(ctx, id)
	if err != nil || caf.StaffID != requestUser.StaffID {
		return globalModels.ErrForbidden
	}
	err = s.surveyRepo.SetSurveyTemplate(ctx, survey, id, requestUser.StaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return globalModels.ErrCafeIsNotExist
		}
		return globalModels.ErrBadJSON
	}
	return err
}

func (s SurveyUsecase) GetSurveyTemplate(ctx context.Context, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	requestUser, err := s.staffUcase.GetFromSession(ctx)
	if err != nil || !requestUser.IsOwner {
		return "", globalModels.ErrForbidden
	}
	survTemplate, err := s.surveyRepo.GetSurveyTemplate(ctx, id, requestUser.StaffID)
	return survTemplate, err
}

func (s SurveyUsecase) SubmitSurvey(ctx context.Context, surv string, customerUUID string) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	err := s.surveyRepo.SubmitSurvey(ctx, surv, customerUUID)
	if err != nil {
		badUUIDmessage := fmt.Sprintf(`pq: invalid input syntax for type uuid: "%s"`, customerUUID)
		if err.Error() == badUUIDmessage {
			return globalModels.ErrBadUuid
		}
		return globalModels.ErrBadJSON
	}
	return err
}