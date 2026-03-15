package httpx

import (
	"encoding/json"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) error

func Adapt(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		if err == nil {
			return
		}

		WriteError(w, err)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	restErr, ok := err.(RestErr)

	if !ok {
		restErr = RestErr{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(restErr.Status)

	json.NewEncoder(w).Encode(restErr)
}

func JSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
