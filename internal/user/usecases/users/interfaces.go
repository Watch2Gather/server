package users

import (
	"context"

	"github.com/google/uuid"

	"github.com/Watch2Gather/server/internal/user/domain"
)

type (
	UserRepo interface {
		Create(context.Context, *domain.RegisterModel) (*domain.User, error)
		Update(context.Context, *domain.User) (*domain.User, error)
		UpdatePassword(context.Context, *domain.ChangePasswordModel) error
		CheckPassword(context.Context, *domain.LoginModel) (uuid.UUID, string, error)
		FindByID(context.Context, uuid.UUID) (*domain.User, error)
		FindByName(context.Context, string) (*domain.User, error)
		UpdateToken(context.Context, *domain.UpdateTokenModel) error
		GetAllFriends(context.Context, uuid.UUID) ([]*domain.User, error)
		AddFriend(context.Context, *domain.AddFriendModel) error
	}
	UseCase interface {
		Login(context.Context, *domain.LoginModel) (*domain.Token, error)
		Register(context.Context, *domain.RegisterModel) error
		ChangePassword(context.Context, *domain.ChangePasswordModel) error
		ChangeUserData(context.Context, *domain.User) (*domain.User, error)
		RefreshToken(context.Context, *domain.Token) (*domain.Token, error)
		GetUserData(context.Context, uuid.UUID) (*domain.UserInfo, error)
		GetAvatar(context.Context, string) (*[]byte, error)
		GetAllFriends(context.Context, uuid.UUID) ([]*domain.User, error)
		AddFriend(context.Context, *domain.AddFriendModel) error
	}
)
