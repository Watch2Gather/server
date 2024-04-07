package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUnauthorized      = errors.New("invalid credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("password is invalid")
)
