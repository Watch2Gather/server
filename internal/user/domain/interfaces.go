package domain

import (
	"context"

	"github.com/google/uuid"
)

type (
	UserRepo interface {
		GetByIDs(context.Context, uuid.UUID) (error)
	}
)
