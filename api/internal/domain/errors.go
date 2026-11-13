package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrDeviceRevoked     = errors.New("device access has been revoked")
	ErrConflict          = errors.New("version conflict detected")
	ErrValidation        = errors.New("validation failed")
)
