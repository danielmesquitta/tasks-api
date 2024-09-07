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
	queries *Queries
}

func NewMySQLTaskRepo(queries *Queries) *MySQLTaskRepo {
	return &MySQLTaskRepo{
		queries: queries,
	}
}

func (m MySQLTaskRepo) GetTaskByID(
	ctx context.Context,
	id string,
) (entity.Task, error) {
	result, err := m.queries.GetTaskByID(ctx, id)

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
		results, err = m.queries.ListTasks(ctx)
	} else {
		results, err = m.queries.ListTasksWithFilters(ctx, sql.NullString{
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

	db := m.queries.getDBorTX(ctx)
	if err := db.CreateTask(ctx, args); err != nil {
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

	db := m.queries.getDBorTX(ctx)
	if err := db.UpdateTask(ctx, args); err != nil {
		return entity.NewErr(err)
	}

	return nil
}

func (m MySQLTaskRepo) DeleteTask(ctx context.Context, id string) error {
	db := m.queries.getDBorTX(ctx)
	if err := db.DeleteTask(ctx, id); err != nil {
		return entity.NewErr(err)
	}

	return nil
}

var _ repo.TaskRepo = (*MySQLTaskRepo)(nil)
