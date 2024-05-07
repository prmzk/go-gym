package auth

import "errors"

var (
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrInvalidEmailOrName  = errors.New("invalid email address and/or name")
	ErrEmailNotFound       = errors.New("email not registered")
	ErrDuplicateEmail      = errors.New("email alredy registered")
	ErrInvalidBearerToken  = errors.New("invalid bearer token")
	ErrAccessTokenNotFound = errors.New("access token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
