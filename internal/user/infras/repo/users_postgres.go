package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/Watch2Gather/server/internal/user/domain"
	"github.com/Watch2Gather/server/internal/user/infras/postgresql"
	"github.com/Watch2Gather/server/internal/user/usecases/users"
	"github.com/Watch2Gather/server/pkg/postgres"
)

const _defaultEntityCap = 64

type userRepo struct {
	pg postgres.DBEngine
}

var _ users.UserRepo = (*userRepo)(nil)

var RepositorySet = wire.NewSet(NewUserRepo)

func NewUserRepo(pg postgres.DBEngine) users.UserRepo {
	return &userRepo{pg: pg}
}

func (u *userRepo) Create(ctx context.Context, model *domain.RegisterModel) (*domain.User, error) {
	querier := postgresql.New(u.pg.GetDB())

	hash, err := bcrypt.GenerateFromPassword([]byte(model.Password), 10)
	if err != nil {
		return nil, errors.Wrap(err, "bcrypt.GenerateFromPassword")
	}

	results, err := querier.CreateUser(ctx, postgresql.CreateUserParams{
		Username: model.Username,
		Email:    model.Email,
		PwdHash:  string(hash),
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.CreateUser")
	}

	user := &domain.User{
		Username: results.Username,
		Email:    results.Email,
	}

	return user, nil
}

func (u *userRepo) CheckPassword(ctx context.Context, model *domain.LoginModel) (uuid.UUID, string, error) {
	querier := postgresql.New(u.pg.GetDB())

	user, err := querier.GetUserByUsername(ctx, model.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", domain.ErrUnauthorized
		}
		return uuid.Nil, "", errors.Wrap(err, "querier.GetUserByUsername")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PwdHash), []byte(model.Password))
	if err != nil {
		return uuid.Nil, "", errors.Wrap(err, "bcrypt.CompareHashAndPassword")
	}

	return user.ID, user.Email, nil
}

func (u *userRepo) Update(ctx context.Context, model *domain.User) (*domain.User, error) {
	querier := postgresql.New(u.pg.GetDB())

	var newUsername, NewEmail, NewAvatar sql.NullString
	if model.Username != "" {
		newUsername = sql.NullString{
			String: model.Username,
			Valid:  true,
		}
	}

	if model.Email != "" {
		NewEmail = sql.NullString{
			String: model.Email,
			Valid:  true,
		}
	}

	if model.Avatar != "" {
		NewEmail = sql.NullString{
			String: model.Avatar,
			Valid:  true,
		}
	}

	updatedUser, err := querier.UpdateUser(ctx, postgresql.UpdateUserParams{
		ID:       model.ID,
		Username: newUsername,
		Email:    NewEmail,
		Avatar:   NewAvatar,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.UpdateUser")
	}

	return &domain.User{
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Avatar:   updatedUser.Avatar.String,
	}, nil
}

func (u *userRepo) UpdatePassword(ctx context.Context, model *domain.ChangePasswordModel) error {
	querier := postgresql.New(u.pg.GetDB())

	pwdHash, err := querier.GetPasswordById(ctx, model.ID)
	if err != nil {
		return errors.Wrap(err, "querier.GetPasswordById")
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(model.OldPasswordHash))
	if err != nil {
		return errors.Wrap(err, "bcrypt.CompareHashAndPassword")
	}

	err = querier.UpdateUserPassword(ctx, postgresql.UpdateUserPasswordParams{
		ID:      model.ID,
		PwdHash: model.NewPasswordHash,
	})
	if err != nil {
		return errors.Wrap(err, "querier.UpdateUserPassword")
	}

	return nil
}

func (u *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	querier := postgresql.New(u.pg.GetDB())

	user, err := querier.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetUserByToken")
	}

	return &domain.User{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar.String,
		ID:       user.ID,
		Token:    user.Token.String,
	}, nil
}

func (u *userRepo) UpdateToken(ctx context.Context, model *domain.UpdateTokenModel) error {
	querier := postgresql.New(u.pg.GetDB())

	str := sql.NullString{
		String: model.RefreshToken,
		Valid:  true,
	}
	err := querier.UpdateToken(ctx, postgresql.UpdateTokenParams{
		ID:    model.ID,
		Token: str,
	})
	if err != nil {
		return errors.Wrap(err, "querier.UpdateToken")
	}

	return nil
}
