package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/msgbroker"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
	"github.com/jinzhu/copier"
)

type FinishTask struct {
	validator *validator.Validator
	msgBroker msgbroker.MessageBroker
	taskRepo  repo.TaskRepo
	userRepo  repo.UserRepo
}

func NewFinishTask(
	validator *validator.Validator,
	msgBroker msgbroker.MessageBroker,
	taskRepo repo.TaskRepo,
	userRepo repo.UserRepo,
) *FinishTask {
	return &FinishTask{
		validator: validator,
		msgBroker: msgBroker,
		taskRepo:  taskRepo,
		userRepo:  userRepo,
	}
}

type FinishTaskParams struct {
	TaskID   string      `json:"task_id,omitempty" validate:"required,uuid"`
	UserID   string      `json:"user_id,omitempty" validate:"required,uuid"`
	UserRole entity.Role `json:"role,omitempty"    validate:"required,min=1,max=2"`
}

func (f *FinishTask) Execute(params FinishTaskParams) error {
	if params.UserRole != entity.RoleTechnician {
		return entity.ErrUserNotAllowedToFinishTask
	}

	if err := f.validator.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return validationErr
	}

	user, err := f.userRepo.GetUserByID(context.Background(), params.UserID)
	if err != nil {
		return entity.NewErr(err)
	}

	if user.ID == "" {
		return entity.ErrUserNotFound
	}

	if user.Role != entity.RoleTechnician {
		return entity.ErrUserNotAllowedToFinishTask
	}

	task, err := f.taskRepo.GetTaskByID(context.Background(), params.TaskID)
	if err != nil {
		return entity.NewErr(err)
	}

	if task.ID == "" {
		return entity.ErrTaskNotFound
	}

	task.FinishedAt = time.Now()

	var repoParams repo.UpdateTaskParams
	if err = copier.Copy(&repoParams, task); err != nil {
		return entity.NewErr(err)
	}

	if err = f.taskRepo.UpdateTask(context.Background(), repoParams); err != nil {
		return entity.NewErr(err)
	}

	task.UpdatedAt = time.Now()

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return entity.NewErr(err)
	}

	if err := f.msgBroker.Publish(msgbroker.TopicTaskFinished, taskBytes); err != nil {
		return entity.NewErr(err)
	}

	return nil
}
