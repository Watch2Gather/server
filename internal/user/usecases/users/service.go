package users

import (
	"context"

	"github.com/google/wire"
	"github.com/pkg/errors"

	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/user/domain"
)

type usecase struct {
	userRepo UserRepo
}

var _ UseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	userRepo UserRepo,
) UseCase {
	return &usecase{
		userRepo: userRepo,
	}
}

func (u *usecase) Login(ctx context.Context, model *domain.LoginModel) (*domain.Token, error) {
	id, err := u.userRepo.CheckPassword(ctx, model)
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "userRepo.CheckPassword")
	}

	var tokens domain.Token

	tokens.AccessToken, err = sharedkernel.CreateAccessToken(ctx, sharedkernel.UserData{
		ID:       id,
		Username: model.Username,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateAccessToken")
	}

	tokens.RefreshToken, err = sharedkernel.CreateRefreshToken(ctx)
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateRefreshToken")
	}

	return &tokens, nil
}

func (u *usecase) Register(ctx context.Context, model *domain.RegisterModel) error {
	if _, err := u.userRepo.Create(ctx, model); err != nil {
		return errors.Wrap(err, "userRepo.Create")
	}

	return nil
}

func (u *usecase) ChangePassword(ctx context.Context, model *domain.ChangePasswordModel) error {
	if err := u.userRepo.UpdatePassword(ctx, model); err != nil {
		return errors.Wrap(err, "userRepo.UpdatePassword")
	}

	return nil
}

func (u *usecase) ChangeUserData(ctx context.Context, model *domain.ChangeUserDataModel) (*domain.User, error) {
	user, err := u.userRepo.Update(ctx, model)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.Update")
	}
	return user, nil
}

func (u *usecase) RefreshToken(ctx context.Context, token *domain.Token) (*domain.Token, error) {
	user, err := u.userRepo.FindByToken(ctx, &domain.RefreshTokenModel{
		RefreshToken: token.RefreshToken,
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.FindByToken")
	}

	var tokens domain.Token

	tokens.RefreshToken, err = sharedkernel.CreateRefreshToken(ctx)
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateRefreshToken")
	}

	err = u.userRepo.UpdateToken(ctx, &domain.UpdateTokenModel{
		ID:           user.ID,
		RefreshToken: token.RefreshToken,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.UpdateToken")
	}

	tokens.AccessToken, err = sharedkernel.CreateAccessToken(ctx, sharedkernel.UserData{
		ID:       user.ID,
		Username: user.Username,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateAccessToken")
	}

	return &tokens, nil
}
