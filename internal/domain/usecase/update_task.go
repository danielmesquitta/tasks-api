package usecase

import (
	"context"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/symcrypt"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type UpdateTask struct {
	validator validator.Validator
	symCrypto symcrypt.SymmetricalEncrypter
	taskRepo  repo.TaskRepo
	userRepo  repo.UserRepo
}

func NewUpdateTask(
	validator validator.Validator,
	symCrypto symcrypt.SymmetricalEncrypter,
	taskRepo repo.TaskRepo,
	userRepo repo.UserRepo,
) *UpdateTask {
	return &UpdateTask{
		validator: validator,
		symCrypto: symCrypto,
		taskRepo:  taskRepo,
		userRepo:  userRepo,
	}
}

type UpdateTaskParams struct {
	ID               string      `json:"id,omitempty"                  validate:"required,uuid"`
	UserID           string      `json:"user_id,omitempty"             validate:"required,uuid"`
	UserRole         entity.Role `json:"user_role,omitempty"           validate:"required,min=1,max=2"`
	Summary          string      `json:"summary,omitempty"             validate:"required,min=1,max=2500"`
	AssignedToUserID *string     `json:"assigned_to_user_id,omitempty" validate:"omitempty,uuid"`
}

func (u *UpdateTask) Execute(params UpdateTaskParams) error {
	if err := u.validator.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return validationErr
	}

	task, err := u.taskRepo.GetTaskByID(context.Background(), params.ID)
	if err != nil {
		return entity.NewErr(err)
	}

	if task.ID == "" {
		return entity.ErrTaskNotFound
	}

	switch params.UserRole {
	case entity.RoleTechnician:
		if params.AssignedToUserID != nil {
			return entity.ErrUserNotAllowedToUpdateAssignedUser
		}

		if task.AssignedToUserID == nil ||
			*task.AssignedToUserID != params.UserID {
			return entity.ErrUserNotAllowedToUpdateTask
		}

	case entity.RoleManager:
		if params.AssignedToUserID == nil {
			break
		}

		var assignedUser entity.User
		assignedUser, err = u.userRepo.GetUserByID(
			context.Background(),
			*params.AssignedToUserID,
		)
		if err != nil {
			return entity.NewErr(err)
		}

		if assignedUser.ID == "" {
			return entity.ErrUserNotFound
		}
	}

	repoParams := repo.UpdateTaskParams{}
	if err = copier.CopyWithOption(&repoParams, params, copier.Option{
		IgnoreEmpty: true,
	}); err != nil {
		return entity.NewErr(err)
	}

	encryptedSummary, err := u.symCrypto.Encrypt(params.Summary)
	if err != nil {
		return entity.NewErr(err)
	}

	repoParams.Summary = encryptedSummary

	if err := u.taskRepo.UpdateTask(context.Background(), repoParams); err != nil {
		return entity.NewErr(err)
	}

	return nil
}
