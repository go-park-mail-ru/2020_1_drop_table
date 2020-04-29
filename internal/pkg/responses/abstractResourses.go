package responses

//easyjson:json
type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//easyjson:json
type HttpResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []HttpError `json:"errors"`
}
