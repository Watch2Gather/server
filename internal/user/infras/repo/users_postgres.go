package repo

import (
	"context"
	"database/sql"

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

func (u *userRepo) Create(ctx context.Context, register *domain.RegisterModel) (*domain.User, error) {
	querier := postgresql.New(u.pg.GetDB())

	hash, err := bcrypt.GenerateFromPassword([]byte(register.Password), 10)
	if err != nil {
		return nil, errors.Wrap(err, "bcrypt.GenerateFromPassword")
	}

	results, err := querier.CreateUser(ctx, postgresql.CreateUserParams{
		Username: register.Username,
		Email:    register.Email,
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

func (u *userRepo) CheckPassword(ctx context.Context, login *domain.LoginModel) error {
	querier := postgresql.New(u.pg.GetDB())

	user, err := querier.GetUserByUsername(ctx, login.Username)
	if err != nil {
		return errors.Wrap(err, "querier.GetUserByUsername")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PwdHash), []byte(login.PasswordHash))
	if err != nil {
		return errors.Wrap(err, "bcrypt.CompareHashAndPassword")
	}

	return nil
}

func (u *userRepo) Update(ctx context.Context, userData *domain.ChangeUserDataModel) (*domain.User, error) {
	querier := postgresql.New(u.pg.GetDB())

	var newUsername, NewEmail, NewAvatar sql.NullString
	if userData.Username != "" {
		newUsername = sql.NullString{
			String: userData.Username,
			Valid:  true,
		}
	}

	if userData.Email != "" {
		NewEmail = sql.NullString{
			String: userData.Email,
			Valid:  true,
		}
	}

	if userData.Avatar != "" {
		NewEmail = sql.NullString{
			String: userData.Avatar,
			Valid:  true,
		}
	}

	updatedUser, err := querier.UpdateUser(ctx, postgresql.UpdateUserParams{
		ID:       userData.ID,
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

func (u *userRepo) UpdatePassword(ctx context.Context, passwordData *domain.ChangePasswordModel) error {
	querier := postgresql.New(u.pg.GetDB())

	pwdHash, err := querier.GetPasswordById(ctx, passwordData.ID)
	if err != nil {
		return errors.Wrap(err, "querier.GetPasswordById")
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(passwordData.OldPasswordHash))
	if err != nil {
		return errors.Wrap(err, "bcrypt.CompareHashAndPassword")
	}

	err = querier.UpdateUserPassword(ctx, postgresql.UpdateUserPasswordParams{
		ID:      passwordData.ID,
		PwdHash: passwordData.NewPasswordHash,
	})
	if err != nil {
		return errors.Wrap(err, "querier.UpdateUserPassword")
	}

	return nil
}
