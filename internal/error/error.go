package error

import (
	"errors"
)

var (
	InternalError = errors.New("internal error")

	// Domain error
	UserAlreadyExists   = errors.New("user already exists")
	UserValidationError = errors.New("validation error")
)

type BaseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (r BaseError) Error() string {
	return r.Message
}

func (r BaseError) Code() int {
	return r.Status
}
