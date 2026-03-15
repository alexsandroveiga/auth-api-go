package httpx

import "net/http"

type RestErr struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func (e RestErr) Error() string {
	return e.Message
}

func BadRequest(message string) RestErr {
	return RestErr{Status: http.StatusBadRequest, Message: message}
}

func UnprocessableEntity(message string) RestErr {
	return RestErr{Status: http.StatusUnprocessableEntity, Message: message}
}

func Unauthorized(message string) RestErr {
	return RestErr{Status: http.StatusUnauthorized, Message: message}
}
