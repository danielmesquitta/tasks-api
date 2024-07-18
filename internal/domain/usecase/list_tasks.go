package usecase

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/pkg/cryptoutil"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
)

type ListTasks struct {
	validator *validator.Validator
	crypto    *cryptoutil.AESCrypto
	taskRepo  repo.TaskRepo
}

func NewListTasks(
	validator *validator.Validator,
	crypto *cryptoutil.AESCrypto,
	taskRepo repo.TaskRepo,
) *ListTasks {
	return &ListTasks{
		validator: validator,
		crypto:    crypto,
		taskRepo:  taskRepo,
	}
}

type ListTasksParams struct {
	UserRole entity.Role `json:"user_role,omitempty" validate:"required,min=1,max=2"`
	UserID   string      `json:"user_id,omitempty"   validate:"required,uuid"`
}

func (l *ListTasks) Execute(
	params ListTasksParams,
) ([]entity.Task, error) {
	if err := l.validator.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return nil, validationErr
	}

	var results []entity.Task
	var err error
	switch params.UserRole {
	case entity.RoleManager:
		results, err = l.taskRepo.ListTasks(context.Background())

	case entity.RoleTechnician:
		results, err = l.taskRepo.ListTasks(
			context.Background(),
			repo.WithAssignedToUserID(params.UserID),
		)

	default:
		return nil, entity.ErrValidation
	}

	if err != nil {
		return nil, entity.NewErr(err)
	}

	tasks := make([]entity.Task, len(results))
	copy(tasks, results)

	for i, task := range tasks {
		decryptedSummary, err := l.crypto.Decrypt(task.Summary)
		if err != nil {
			return nil, entity.NewErr(err)
		}
		tasks[i].Summary = decryptedSummary
	}

	return tasks, nil
}
