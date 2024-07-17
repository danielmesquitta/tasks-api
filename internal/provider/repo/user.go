package repo

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}
