package v1

import (
	"encoding/json"
	"fio/internal/domain"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type idResponse struct {
	ID interface{} `json:"id"`
}

type getPersonsResponse struct {
	Data []domain.Person `json:"data"`
}

func newErrorResponse(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(errorResponse{msg}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}

func newStatusReponse(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(statusResponse{msg}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}

func newGetPersonsResponse(w http.ResponseWriter, persons []domain.Person, status int) {
	resp, _ := json.Marshal(getPersonsResponse{persons}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}
