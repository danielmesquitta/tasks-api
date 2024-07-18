package usecase

import (
	"context"
	"sync"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/pkg/crypto"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
	"github.com/jinzhu/copier"
)

type CreateTask struct {
	validator *validator.Validator
	crypto    *crypto.Crypto
	taskRepo  repo.TaskRepo
	userRepo  repo.UserRepo
}

func NewCreateTask(
	validator *validator.Validator,
	crypto *crypto.Crypto,
	taskRepo repo.TaskRepo,
	userRepo repo.UserRepo,
) *CreateTask {
	return &CreateTask{
		validator: validator,
		crypto:    crypto,
		taskRepo:  taskRepo,
		userRepo:  userRepo,
	}
}

type CreateTaskParams struct {
	UserRole         entity.Role `json:"user_role,omitempty"           validate:"required,min=1,max=2"`
	Summary          string      `json:"summary,omitempty"             validate:"required,max=2500"`
	CreatedByUserID  string      `json:"created_by_user_id,omitempty"  validate:"required,uuid"`
	AssignedToUserID string      `json:"assigned_to_user_id,omitempty" validate:"uuid,omitempty"`
}

func (c *CreateTask) Execute(params CreateTaskParams) (ID string, err error) {
	if params.UserRole != entity.RoleManager {
		return "", entity.ErrUserNotAllowedToCreateTask
	}

	if err = c.validator.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return "", validationErr
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var createdByUser, assignedToUser entity.User
	var userRepoErr error
	go func() {
		defer wg.Done()
		createdByUser, err = c.userRepo.GetUserByID(
			context.Background(),
			params.CreatedByUserID,
		)
		if err != nil {
			userRepoErr = err
		}
	}()

	go func() {
		defer wg.Done()
		assignedToUser, err = c.userRepo.GetUserByID(
			context.Background(),
			params.AssignedToUserID,
		)
		if err != nil {
			userRepoErr = err
		}
	}()

	wg.Wait()

	if userRepoErr != nil {
		return "", entity.NewErr(err)
	}

	if createdByUser.ID == "" {
		return "", entity.ErrCreatedByUserNotFound
	}

	if assignedToUser.ID == "" {
		return "", entity.ErrAssignToUserNotFound
	}

	if createdByUser.Role != entity.RoleManager {
		return "", entity.ErrUserNotAllowedToCreateTask
	}

	if assignedToUser.Role != entity.RoleTechnician {
		return "", entity.ErrInvalidRoleForAssignedUser
	}

	encryptedSummary, err := c.crypto.Encrypt(params.Summary)
	if err != nil {
		return "", entity.NewErr(err)
	}

	params.Summary = encryptedSummary

	repoParams := repo.CreateTaskParams{}
	if err = copier.Copy(&repoParams, params); err != nil {
		return "", entity.NewErr(err)
	}

	taskID, err := c.taskRepo.CreateTask(context.Background(), repoParams)
	if err != nil {
		return "", entity.NewErr(err)
	}

	return taskID, nil
}
