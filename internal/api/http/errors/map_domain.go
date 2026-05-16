package errors

import (
	"auth/internal/api/http/dto"
)

type ErrorBody interface {
	SetRequestID(string)
}

type MappedError struct {
	Status int
	Body   ErrorBody
}

func NewBasicMapped(status int, code, message string) MappedError {
	return MappedError{
		Status: status,
		Body: &dto.ErrorResponse{
			Code:    code,
			Message: message,
		},
	}
}
