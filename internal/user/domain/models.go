package domain

type LoginModel struct {
	Username     string
	PasswordHash string
}

type RegisterModel struct {
	Username     string
	PasswordHash string
	Email        string
}

type ChangePasswordModel struct {
	OldPasswordHash string
	NewPasswordHash string
}
