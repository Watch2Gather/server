package usecases

import (
	"context"

	"github.com/Watch2Gather/server/internal/user/domain"
)

type (
	UseCase interface {
		Login(context.Context, *domain.LoginModel) (*domain.Token, error)
		Register(context.Context, *domain.RegisterModel) error
		ChangePassword(context.Context, *domain.ChangePasswordModel) error
	}
)
