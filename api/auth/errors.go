package auth

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrDuplicateEmail     = errors.New("email alredy registered")
	ErrInvalidBearerToken = errors.New("invalid bearer token")
)
