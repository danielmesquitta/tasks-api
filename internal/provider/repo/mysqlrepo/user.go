package mysqlrepo

import (
	"context"
	"database/sql"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/db/mysqldb"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type MySQLUserRepo struct {
	queries *Queries
}

func NewMySQLUserRepo(queries *Queries) *MySQLUserRepo {
	return &MySQLUserRepo{
		queries: queries,
	}
}

func (m MySQLUserRepo) GetUserByID(
	ctx context.Context,
	id string,
) (entity.User, error) {
	result, err := m.queries.GetUserByID(ctx, id)

	if err == sql.ErrNoRows {
		return entity.User{}, nil
	}

	if err != nil {
		return entity.User{}, entity.NewErr(err)
	}

	user := entity.User{}
	if err := copier.Copy(&user, result); err != nil {
		return entity.User{}, entity.NewErr(err)
	}

	return user, nil
}

func (m MySQLUserRepo) GetUserByEmail(
	ctx context.Context,
	email string,
) (entity.User, error) {
	result, err := m.queries.GetUserByEmail(ctx, email)

	if err == sql.ErrNoRows {
		return entity.User{}, nil
	}

	if err != nil {
		return entity.User{}, entity.NewErr(err)
	}

	user := entity.User{}
	if err := copier.Copy(&user, result); err != nil {
		return entity.User{}, entity.NewErr(err)
	}

	return user, nil
}

func (m MySQLUserRepo) CreateUser(
	ctx context.Context,
	params repo.CreateUserParams,
) error {
	args := mysqldb.CreateUserParams{}
	if err := copier.Copy(&args, params); err != nil {
		return entity.NewErr(err)
	}

	db := m.queries.getDBorTX(ctx)
	if err := db.CreateUser(ctx, args); err != nil {
		return entity.NewErr(err)
	}

	return nil
}

var _ repo.UserRepo = (*MySQLUserRepo)(nil)
