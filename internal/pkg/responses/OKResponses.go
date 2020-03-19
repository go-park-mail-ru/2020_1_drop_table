package responses

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

func SendOKAnswer(data interface{}, w http.ResponseWriter) {
	serializedData, err := json.Marshal(HttpResponse{Data: data})
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
	log.Info().Msgf("OK message sent")
}
