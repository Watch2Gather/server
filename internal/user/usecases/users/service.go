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

	var tokens *domain.Token

	tokens.AccessToken, err = sharedkernel.CreateAccessToken(sharedkernel.UserData{
		ID:       id,
		Username: model.Username,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateAccessToken")
	}

	tokens.RefreshToken, err = sharedkernel.CreateRefreshToken()
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateRefreshToken")
	}

	return tokens, nil
}

func (u *usecase) Register(_ context.Context, _ *domain.RegisterModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (u *usecase) ChangePassword(_ context.Context, _ *domain.ChangePasswordModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (u *usecase) ChangeUserData(_ context.Context, _ *domain.ChangeUserDataModel) (_ error) {
	panic("not implemented") // TODO: Implement
}

func (u *usecase) RefreshToken(context.Context, *dom)
