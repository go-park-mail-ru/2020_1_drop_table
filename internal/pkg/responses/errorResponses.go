package responses

import (
	globalModels "2020_1_drop_table/internal/app/models"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (e *HttpError) Error() string {
	return fmt.Sprintf("Error: '%s', with status code: %d", e.Message, e.Code)
}

func SendServerError(errorMessage string, w http.ResponseWriter) {
	log.Error().Msgf(errorMessage)
	w.WriteHeader(http.StatusInternalServerError)
}

func SendSingleError(errorMessage string, w http.ResponseWriter) {
	log.Error().Msgf(errorMessage)
	errs := make([]HttpError, 1, 1)
	errs[0] = HttpError{
		Code:    400,
		Message: errorMessage,
	}
	SendSeveralErrors(errs, w)
}

func SendSingleErrorWithMessage(errForLog error, errorMessage string, w http.ResponseWriter) {
	log.Error().Msgf(errForLog.Error())
	errs := make([]HttpError, 1, 1)
	errs[0] = HttpError{
		Code:    400,
		Message: errorMessage,
	}
	SendSeveralErrors(errs, w)
}

func SendSeveralErrors(errors []HttpError, w http.ResponseWriter) {
	httpResponse := HttpResponse{Errors: errors}
	serializedError, err := httpResponse.MarshalJSON()
	if err != nil {
		message := fmt.Sprintf("HttpResponse is json serializing: %s", err.Error())
		SendServerError(message, w)
		return
	}

	_, err = w.Write(serializedError)
	if err != nil {
		message := fmt.Sprintf("HttpResponse while writing is socket: %s", err.Error())
		SendServerError(message, w)
		return
	}
}

func SendForbidden(w http.ResponseWriter) {
	SendSingleError(globalModels.ErrForbidden.Error(), w)
}
