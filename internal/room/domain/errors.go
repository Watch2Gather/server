package domain

import "errors"

var (
	ErrInvalidID  = errors.New("invalid uuid format")
	ErrNoRoomOpen = errors.New("no open room to send message")
)
