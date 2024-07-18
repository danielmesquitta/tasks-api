package mysqlrepo

import (
	"context"
	"database/sql"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/db/mysqldb"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type MySQLTaskRepo struct {
	db *mysqldb.Queries
}

func NewMySQLTaskRepo(db *mysqldb.Queries) *MySQLTaskRepo {
	return &MySQLTaskRepo{
		db: db,
	}
}

func (m MySQLTaskRepo) GetTaskByID(
	ctx context.Context,
	id string,
) (entity.Task, error) {
	result, err := m.db.GetTaskByID(ctx, id)

	if err == sql.ErrNoRows {
		return entity.Task{}, nil
	}

	if err != nil {
		return entity.Task{}, entity.NewErr(err)
	}

	task := entity.Task{}
	if err := copier.Copy(&task, result); err != nil {
		return entity.Task{}, entity.NewErr(err)
	}

	return task, nil
}

func (m MySQLTaskRepo) ListTasks(
	ctx context.Context,
	opts ...repo.ListTasksOption,
) (tasks []entity.Task, err error) {
	params := repo.ListTasksParams{}
	for _, opt := range opts {
		opt(&params)
	}

	var results []mysqldb.Task
	if params.AssignedToUserID == "" {
		results, err = m.db.ListTasks(ctx)
	} else {
		results, err = m.db.ListTasksWithFilters(ctx, sql.NullString{
			String: params.AssignedToUserID,
			Valid:  true,
		})
	}

	if err != nil {
		return nil, entity.NewErr(err)
	}

	for _, result := range results {
		task := entity.Task{}
		if err := copier.Copy(&task, result); err != nil {
			return nil, entity.NewErr(err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (m MySQLTaskRepo) CreateTask(
	ctx context.Context,
	params repo.CreateTaskParams,
) error {
	args := mysqldb.CreateTaskParams{}
	if err := copier.Copy(&args, params); err != nil {
		return entity.NewErr(err)
	}

	if err := m.db.CreateTask(ctx, args); err != nil {
		return entity.NewErr(err)
	}

	return nil
}

func (m MySQLTaskRepo) UpdateTask(
	ctx context.Context,
	params repo.UpdateTaskParams,
) error {
	args := mysqldb.UpdateTaskParams{}
	if err := copier.Copy(&args, params); err != nil {
		return entity.NewErr(err)
	}

	if err := m.db.UpdateTask(ctx, args); err != nil {
		return entity.NewErr(err)
	}

	return nil
}

func (m MySQLTaskRepo) DeleteTask(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}
