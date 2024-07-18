package repo

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
)

type CreateUserParams struct {
	Name     string
	Role     entity.Role
	Email    string
	Password string
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(
		ctx context.Context,
		params CreateUserParams,
	) error
}
