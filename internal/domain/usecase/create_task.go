package usecase

import (
	"context"
	"errors"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/symcrypt"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type CreateTask struct {
	validator validator.Validator
	symCrypto symcrypt.SymmetricalEncrypter
	taskRepo  repo.TaskRepo
	userRepo  repo.UserRepo
}

func NewCreateTask(
	validator validator.Validator,
	symCrypto symcrypt.SymmetricalEncrypter,
	taskRepo repo.TaskRepo,
	userRepo repo.UserRepo,
) *CreateTask {
	return &CreateTask{
		validator: validator,
		symCrypto: symCrypto,
		taskRepo:  taskRepo,
		userRepo:  userRepo,
	}
}

type CreateTaskParams struct {
	UserRole         entity.Role `json:"user_role,omitempty"           validate:"required,min=1,max=2"`
	Summary          string      `json:"summary,omitempty"             validate:"required,max=2500"`
	CreatedByUserID  string      `json:"created_by_user_id,omitempty"  validate:"required,uuid"`
	AssignedToUserID string      `json:"assigned_to_user_id,omitempty" validate:"omitempty,uuid"`
}

func (c *CreateTask) Execute(
	ctx context.Context,
	params CreateTaskParams,
) error {
	if params.UserRole != entity.RoleManager {
		return entity.ErrUserNotAllowedToCreateTask
	}

	if err := c.validator.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return validationErr
	}

	var createdByUser, assignedToUser entity.User
	errCh := make(chan error)

	go func() {
		var err error
		defer func() {
			errCh <- err
		}()
		createdByUser, err = c.userRepo.GetUserByID(
			ctx,
			params.CreatedByUserID,
		)
	}()

	assignedUserIsDefined := params.AssignedToUserID != ""
	go func() {
		var err error
		defer func() {
			errCh <- err
		}()
		if !assignedUserIsDefined {
			return
		}
		assignedToUser, err = c.userRepo.GetUserByID(
			ctx,
			params.AssignedToUserID,
		)
	}()

	var errs []error
	routinesCount := 2
	for i := 0; i < routinesCount; i++ {
		err := <-errCh
		if err != nil {
			errs = append(errs, err)
		}
	}

	close(errCh)

	if len(errs) > 0 {
		return entity.NewErr(errors.Join(errs...))
	}

	if createdByUser.ID == "" {
		return entity.ErrCreatedByUserNotFound
	}

	if createdByUser.Role != entity.RoleManager {
		return entity.ErrUserNotAllowedToCreateTask
	}

	if assignedUserIsDefined {
		if assignedToUser.ID == "" {
			return entity.ErrAssignToUserNotFound
		}

		if assignedToUser.Role != entity.RoleTechnician {
			return entity.ErrInvalidRoleForAssignedUser
		}
	}

	encryptedSummary, err := c.symCrypto.Encrypt(params.Summary)
	if err != nil {
		return entity.NewErr(err)
	}

	params.Summary = encryptedSummary

	repoParams := repo.CreateTaskParams{}
	if err = copier.CopyWithOption(&repoParams, params, copier.Option{
		IgnoreEmpty: true,
	}); err != nil {
		return entity.NewErr(err)
	}

	if err := c.taskRepo.CreateTask(ctx, repoParams); err != nil {
		return entity.NewErr(err)
	}

	return nil
}
