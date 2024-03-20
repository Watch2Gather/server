package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("invalid credentials")
