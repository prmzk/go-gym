package auth

import "errors"

var (
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrEmailNotFound       = errors.New("email not registered")
	ErrDuplicateEmail      = errors.New("email alredy registered")
	ErrInvalidBearerToken  = errors.New("invalid bearer token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
