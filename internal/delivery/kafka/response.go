package kafka

import "fio/internal/domain"

type errorResponse struct {
	Message string `json:"message"`
}

type personErrorResponse struct {
	domain.Person
	errorResponse
}
