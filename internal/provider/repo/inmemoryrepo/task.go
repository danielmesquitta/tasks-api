package inmemoryrepo

import (
	"context"
	"slices"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type InMemoryTaskRepo struct {
	Tasks []entity.Task
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{
		Tasks: []entity.Task{},
	}
}

func (im *InMemoryTaskRepo) GetTaskByID(
	_ context.Context,
	id string,
) (entity.Task, error) {
	for _, task := range im.Tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return entity.Task{}, nil
}

func (im *InMemoryTaskRepo) ListTasks(
	_ context.Context,
	opts ...repo.ListTasksOption,
) ([]entity.Task, error) {
	params := repo.ListTasksParams{}
	for _, opt := range opts {
		opt(&params)
	}

	if params.AssignedToUserID != "" {
		var tasks []entity.Task
		for _, task := range im.Tasks {
			if task.AssignedToUserID != nil &&
				*task.AssignedToUserID == params.AssignedToUserID {
				tasks = append(tasks, task)
			}
		}
		return tasks, nil
	}

	return im.Tasks, nil
}

func (im *InMemoryTaskRepo) CreateTask(
	_ context.Context,
	params repo.CreateTaskParams,
) error {
	task := entity.Task{}
	if err := copier.Copy(&task, params); err != nil {
		return entity.NewErr(err)
	}

	task.ID = uuid.NewString()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	im.Tasks = append(im.Tasks, task)

	return nil
}

func (im *InMemoryTaskRepo) UpdateTask(
	_ context.Context,
	params repo.UpdateTaskParams,
) error {
	for i, task := range im.Tasks {
		if task.ID != params.ID {
			continue
		}

		if err := copier.CopyWithOption(
			&task,
			params,
			copier.Option{IgnoreEmpty: true},
		); err != nil {
			return entity.NewErr(err)
		}

		task.UpdatedAt = time.Now()

		im.Tasks[i] = task
		break
	}

	return nil
}

func (im *InMemoryTaskRepo) DeleteTask(_ context.Context, id string) error {
	for i, task := range im.Tasks {
		if task.ID == id {
			im.Tasks = slices.Delete(im.Tasks, i, i+1)
			break
		}
	}

	return nil
}
