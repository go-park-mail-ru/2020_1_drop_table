package responses

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func SendOKAnswer(data interface{}, w http.ResponseWriter) {
	serializedData, err := HttpResponse{Data: data}.MarshalJSON()
	if err != nil {
		log.Error().Msgf(err.Error())
		SendServerError("Server JSON encoding error", w)
	}
	_, err = w.Write(serializedData)
	if err != nil {
		message := fmt.Sprintf("HttpResponse while writing is socket: %s", err.Error())
		SendServerError(message, w)
		return
	}
}
