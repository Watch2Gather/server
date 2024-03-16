package users

import (
	"context"

	"github.com/google/wire"
	"github.com/pkg/errors"

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

func (u *usecase) Login(ctx context.Context, model *domain.LoginModel) (_ *domain.Token, _ error) {
	err := u.userRepo.CheckPassword(ctx, model)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.CheckPassword")
	}
	panic("not implemented") // TODO: Implement
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
