package repo

import (
	"context"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
)

type CreateTaskParams struct {
	Summary          string  `json:"summary"`
	CreatedByUserID  string  `json:"created_by_user_id"`
	AssignedToUserID *string `json:"assigned_to_user_id"`
}

type UpdateTaskParams struct {
	ID               string    `json:"id"`
	Summary          string    `json:"summary"`
	AssignedToUserID *string   `json:"assigned_to_user_id"`
	FinishedAt       time.Time `json:"finished_at"`
}

type ListTasksParams struct {
	AssignedToUserID string `json:"assigned_to_user_id"`
}

type ListTasksOption func(*ListTasksParams)

func WithAssignedToUserID(assignedToUserId string) ListTasksOption {
	return func(params *ListTasksParams) {
		params.AssignedToUserID = assignedToUserId
	}
}

type TaskRepo interface {
	GetTaskByID(ctx context.Context, id string) (entity.Task, error)
	ListTasks(
		ctx context.Context,
		opts ...ListTasksOption,
	) ([]entity.Task, error)
	CreateTask(
		ctx context.Context,
		params CreateTaskParams,
	) error
	UpdateTask(
		ctx context.Context,
		params UpdateTaskParams,
	) error
	DeleteTask(ctx context.Context, id string) error
}
