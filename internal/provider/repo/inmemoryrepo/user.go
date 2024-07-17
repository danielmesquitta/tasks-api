package inmemoryrepo

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
)

type InMemoryUserRepo struct {
	Users []entity.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		Users: []entity.User{},
	}
}

func (im *InMemoryUserRepo) GetUserByID(
	ctx context.Context,
	id string,
) (entity.User, error) {
	for _, user := range im.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return entity.User{}, nil
}
