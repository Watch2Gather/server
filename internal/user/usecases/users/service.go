package users

import (
	"context"

	"github.com/google/uuid"
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
	id, email, err := u.userRepo.CheckPassword(ctx, model)
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "userRepo.CheckPassword")
	}

	var tokens domain.Token

	tokens.AccessToken, err = sharedkernel.CreateAccessToken(ctx, sharedkernel.UserData{
		ID:       id,
		Username: model.Username,
		Email:    email,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateAccessToken")
	}

	tokens.RefreshToken, err = sharedkernel.CreateRefreshToken(ctx, id)
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateRefreshToken")
	}
	err = u.userRepo.UpdateToken(ctx, &domain.UpdateTokenModel{
		ID:           id,
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.UpdateToken")
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

func (u *usecase) ChangeUserData(ctx context.Context, model *domain.User) (*domain.User, error) {
	user, err := u.userRepo.Update(ctx, model)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.Update")
	}
	return user, nil
}

func (u *usecase) RefreshToken(ctx context.Context, model *domain.Token) (*domain.Token, error) {
	tokenString, err := sharedkernel.GetToken(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.GetToken")
	}

	claims, err := sharedkernel.ParseToken(ctx, tokenString)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.ParseToken")
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse")
	}

	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.FindByToken")
	}

	if user.Token != model.RefreshToken {
		return nil, sharedkernel.ErrInvalidToken
	}

	var tokens domain.Token

	tokens.AccessToken, err = sharedkernel.CreateAccessToken(ctx, sharedkernel.UserData{
		Username: user.Username,
		Email:    user.Email,
		ID:       id,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateAccessToken")
	}

	tokens.RefreshToken, err = sharedkernel.CreateRefreshToken(ctx, id)
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.CreateRefreshToken")
	}

	err = u.userRepo.UpdateToken(ctx, &domain.UpdateTokenModel{
		ID:           id,
		RefreshToken: model.RefreshToken,
	})
	if err != nil {
		return &domain.Token{}, errors.Wrap(err, "sharedkernel.UpdateToken")
	}

	return &tokens, nil
}
