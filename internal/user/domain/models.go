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

type RefreshTokenModel struct {
	RefreshToken string
}

type UpdateTokenModel struct {
	ID           uuid.UUID
	RefreshToken string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

type User struct {
	Username string
	Email    string
	Avatar   string
	ID       uuid.UUID
}
