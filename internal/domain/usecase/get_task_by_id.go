package usecase

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/symcrypt"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
)

type GetTaskByID struct {
	validator validator.Validator
	symCrypto symcrypt.SymmetricalEncrypter
	taskRepo  repo.TaskRepo
}

func NewGetTaskByID(
	validator validator.Validator,
	symCrypto symcrypt.SymmetricalEncrypter,
	taskRepo repo.TaskRepo,
) *GetTaskByID {
	return &GetTaskByID{
		validator: validator,
		symCrypto: symCrypto,
		taskRepo:  taskRepo,
	}
}

type GetTaskByIDParams struct {
	ID       string      `json:"id,omitempty"        validate:"required,uuid"`
	UserID   string      `json:"user_id,omitempty"   validate:"required,uuid"`
	UserRole entity.Role `json:"user_role,omitempty" validate:"required,min=1,max=2"`
}

func (u *GetTaskByID) Execute(params GetTaskByIDParams) (entity.Task, error) {
	if err := u.validator.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return entity.Task{}, validationErr
	}

	task, err := u.taskRepo.GetTaskByID(context.Background(), params.ID)
	if err != nil {
		return entity.Task{}, entity.NewErr(err)
	}

	if task.ID == "" {
		return entity.Task{}, entity.ErrTaskNotFound
	}

	if params.UserRole == entity.RoleTechnician {
		if task.AssignedToUserID == nil {
			return entity.Task{}, entity.ErrUserNotAllowedToViewTask
		}

		if *task.AssignedToUserID != params.UserID {
			return entity.Task{}, entity.ErrUserNotAllowedToViewTask
		}
	}

	decryptedSummary, err := u.symCrypto.Decrypt(task.Summary)
	if err != nil {
		return entity.Task{}, entity.NewErr(err)
	}

	task.Summary = decryptedSummary

	return task, nil
}
