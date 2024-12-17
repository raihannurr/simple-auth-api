package errors

import "errors"

var (
	Join = errors.Join
	New  = errors.New

	ErrInternalServerError     = errors.New("internal server error")
	ErrUserNotFound            = errors.New("user not found")
	ErrDestinationUserNotFound = errors.New("destination user not found")
	ErrUserExists              = errors.New("user already exists")
	ErrInsufficientBalance     = errors.New("insufficient balance")
	ErrTokenExpired            = errors.New("token expired")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidLoginCredentials = errors.New("invalid login credentials")
	ErrInvalidPassword         = errors.New("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character")
)
