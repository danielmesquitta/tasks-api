package usecase

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
)

type DeleteTask struct {
	validator validator.Validator
	taskRepo  repo.TaskRepo
}

func NewDeleteTask(
	validator validator.Validator,
	taskRepo repo.TaskRepo,
) *DeleteTask {
	return &DeleteTask{
		validator: validator,
		taskRepo:  taskRepo,
	}
}

type DeleteTaskParams struct {
	TaskID   string      `json:"task_id,omitempty" validate:"required,uuid"`
	UserRole entity.Role `json:"role,omitempty"    validate:"required,min=1,max=2"`
}

func (d *DeleteTask) Execute(
	ctx context.Context,
	params DeleteTaskParams,
) error {
	if err := d.validator.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return validationErr
	}

	if params.UserRole != entity.RoleManager {
		return entity.ErrUserNotAllowedToDeleteTask
	}

	task, err := d.taskRepo.GetTaskByID(ctx, params.TaskID)
	if err != nil {
		return entity.NewErr(err)
	}

	if task.ID == "" {
		return entity.ErrTaskNotFound
	}

	if err := d.taskRepo.DeleteTask(ctx, task.ID); err != nil {
		return entity.NewErr(err)
	}

	return nil
}
