package error

import "errors"

var (
	InternalError = errors.New("internal error")

	// Domain error
	UserAlreadyExists = errors.New("user already exists")
)
