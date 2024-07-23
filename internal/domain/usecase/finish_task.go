package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/broker"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type FinishTask struct {
	validator validator.Validator
	msgBroker broker.MessageBroker
	taskRepo  repo.TaskRepo
}

func NewFinishTask(
	validator validator.Validator,
	msgBroker broker.MessageBroker,
	taskRepo repo.TaskRepo,
) *FinishTask {
	return &FinishTask{
		validator: validator,
		msgBroker: msgBroker,
		taskRepo:  taskRepo,
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

	task, err := f.taskRepo.GetTaskByID(context.Background(), params.TaskID)
	if err != nil {
		return entity.NewErr(err)
	}

	if task.ID == "" {
		return entity.ErrTaskNotFound
	}

	finishedAt := time.Now()
	task.FinishedAt = &finishedAt

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

	if err := f.msgBroker.Publish(broker.TopicTaskFinished, taskBytes); err != nil {
		return entity.NewErr(err)
	}

	return nil
}
