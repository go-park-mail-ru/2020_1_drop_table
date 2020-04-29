package usecase

import (
	cafeClient "2020_1_drop_table/internal/app/cafe/delivery/grpc/client"
	globalModels "2020_1_drop_table/internal/app/models"
	staffClient "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	"2020_1_drop_table/internal/microservices/survey"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SurveyUsecase struct {
	surveyRepo     survey.Repository
	staffClient    *staffClient.StaffClient
	cafeClient     *cafeClient.CafeGRPC
	contextTimeout time.Duration
}

func NewSurveyUsecase(surveyRepo survey.Repository, staffClient *staffClient.StaffClient, cafeClient *cafeClient.CafeGRPC, timeout time.Duration) survey.Usecase {
	return &SurveyUsecase{
		surveyRepo:     surveyRepo,
		contextTimeout: timeout,
		staffClient:    staffClient,
		cafeClient:     cafeClient,
	}
}

func (s SurveyUsecase) SetSurveyTemplate(ctx context.Context, survey string, id int) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	requestUser, err := s.staffClient.GetFromSession(ctx)
	fmt.Println(requestUser)

	if err != nil || !requestUser.IsOwner {
		return globalModels.ErrForbidden
	}
	caf, err := s.cafeClient.GetByID(ctx, id)
	if err != nil || caf.StaffID != requestUser.StaffID {
		return globalModels.ErrForbidden
	}
	err = s.surveyRepo.SetSurveyTemplate(ctx, survey, id, requestUser.StaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return globalModels.ErrCafeIsNotExist
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "surveytemplate_cafeid_key"` {
			err := s.surveyRepo.UpdateSurveyTemplate(ctx, survey, id)
			if err == nil {
				return nil
			}
		}
		return globalModels.ErrBadJSON
	}
	return err
}

func (s SurveyUsecase) GetSurveyTemplate(ctx context.Context, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	survTemplate, err := s.surveyRepo.GetSurveyTemplate(ctx, id)
	if err == sql.ErrNoRows {
		return "", nil
	}
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
