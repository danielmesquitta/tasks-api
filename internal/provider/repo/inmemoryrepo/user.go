package inmemoryrepo

import (
	"context"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type InMemoryUserRepo struct {
	Users []entity.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		Users: []entity.User{},
	}
}

func (im *InMemoryUserRepo) CreateUser(
	_ context.Context,
	params repo.CreateUserParams,
) error {
	user := entity.User{}
	if err := copier.Copy(&user, params); err != nil {
		return entity.NewErr(err)
	}

	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	im.Users = append(im.Users, user)

	return nil
}

func (im *InMemoryUserRepo) GetUserByID(
	_ context.Context,
	id string,
) (entity.User, error) {
	for _, user := range im.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return entity.User{}, nil
}

func (im *InMemoryUserRepo) GetUserByEmail(
	_ context.Context,
	email string,
) (entity.User, error) {
	for _, user := range im.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return entity.User{}, nil
}

var _ repo.UserRepo = (*InMemoryUserRepo)(nil)
