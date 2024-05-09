package domain

import (
	"context"

	"github.com/google/uuid"
)

type (
	UserInfoDomainService interface {
		GetUserByID(context.Context, uuid.UUID) (*UserModel, error)
	}
)
