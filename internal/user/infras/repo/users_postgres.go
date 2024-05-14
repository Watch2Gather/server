package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/lib/pq"
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

func NewUserRepo(
	pg postgres.DBEngine,
) users.UserRepo {
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, domain.ErrUserAlreadyExists
			}
		}
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
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, domain.ErrUserAlreadyExists
			}
		}
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

	oldPwdHash, err := querier.GetPasswordById(ctx, model.ID)
	if err != nil {
		return errors.Wrap(err, "querier.GetPasswordById")
	}

	err = bcrypt.CompareHashAndPassword([]byte(oldPwdHash), []byte(model.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domain.ErrInvalidPassword
		}
		return errors.Wrap(err, "bcrypt.CompareHashAndPassword")
	}

	newPwdHash, err := bcrypt.GenerateFromPassword([]byte(model.NewPassword), 10)
	if err != nil {
		return errors.Wrap(err, "bcrypt.GenerateFromPassword")
	}
	err = querier.UpdateUserPassword(ctx, postgresql.UpdateUserPasswordParams{
		ID:      model.ID,
		PwdHash: string(newPwdHash),
	})
	if err != nil {
		return errors.Wrap(err, "querier.UpdateUserPassword")
	}

	return nil
}

func (u *userRepo) FindByName(ctx context.Context, name string) (*domain.User, error) {
	querier := postgresql.New(u.pg.GetDB())

	user, err := querier.GetUserByUsername(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetUserByName")
	}

	return &domain.User{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar.String,
		ID:       user.ID,
	}, nil
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

func (u *userRepo) GetAllFriends(ctx context.Context, id uuid.UUID) ([]*domain.User, error) {
	querier := postgresql.New(u.pg.GetDB())

	friends, err := querier.GetFriendList(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetFriendList")
	}

	var res []*domain.User
	for _, friend := range friends {
		res = append(res, &domain.User{
			Username: friend.Username,
			Avatar:   friend.Avatar.String,
			ID:       friend.ID,
		})
	}

	return res, nil
}

func (u *userRepo) AddFriend(ctx context.Context, model *domain.AddFriendModel) error {
	querier := postgresql.New(u.pg.GetDB())

	err := querier.AddFriendById(ctx, postgresql.AddFriendByIdParams{
		UserID1: model.UserID,
		UserID2: model.FriendID,
	})
	if err != nil {
		return errors.Wrap(err, "querier.AddFriendById1")
	}

	err = querier.AddFriendById(ctx, postgresql.AddFriendByIdParams{
		UserID2: model.UserID,
		UserID1: model.FriendID,
	})
	if err != nil {
		return errors.Wrap(err, "querier.AddFriendById2")
	}

	return nil
}
