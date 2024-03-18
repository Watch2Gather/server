package users

import (
	"context"

	"github.com/google/uuid"

	"github.com/Watch2Gather/server/internal/user/domain"
)

type (
	UserRepo interface {
		Create(context.Context, *domain.RegisterModel) (*domain.User, error)
		Update(context.Context, *domain.ChangeUserDataModel) (*domain.User, error)
		UpdatePassword(context.Context, *domain.ChangePasswordModel) error
		CheckPassword(context.Context, *domain.LoginModel) (uuid.UUID, error)
	}
	UseCase interface {
		Login(context.Context, *domain.LoginModel) (*domain.Token, error)
		Register(context.Context, *domain.RegisterModel) error
		ChangePassword(context.Context, *domain.ChangePasswordModel) error
		ChangeUserData(context.Context, *domain.ChangeUserDataModel) error
	}
)