package domain

import "github.com/google/uuid"

type LoginModel struct {
	Username string
	Password string
}

type RegisterModel struct {
	Username string
	Password string
	Email    string
}

type ChangePasswordModel struct {
	OldPassword string
	NewPassword string
	ID          uuid.UUID
}

type User struct {
	Username string
	Email    string
	Avatar   string
	Token    string
	ID       uuid.UUID
}

type UpdateTokenModel struct {
	RefreshToken string
	ID           uuid.UUID
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
