package domain

import "github.com/google/uuid"

type LoginModel struct {
	Username     string
	PasswordHash string
}

type RegisterModel struct {
	Username string
	Password string
	Email    string
}

type ChangePasswordModel struct {
	ID              uuid.UUID
	OldPasswordHash string
	NewPasswordHash string
}

type ChangeUserDataModel struct {
	ID       uuid.UUID
	Username string
	Email    string
	Avatar   string
}

type User struct {
	Username string
	Email    string
	Avatar   string
}
