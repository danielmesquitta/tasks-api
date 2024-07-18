package usecase

import (
	"context"
	"strings"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/pkg/cryptoutil"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
	"github.com/jinzhu/copier"
)

type CreateUser struct {
	val      *validator.Validator
	bcr      *cryptoutil.Bcrypt
	userRepo repo.UserRepo
}

func NewCreateUser(
	val *validator.Validator,
	bcr *cryptoutil.Bcrypt,
	userRepo repo.UserRepo,
) *CreateUser {
	return &CreateUser{
		val:      val,
		bcr:      bcr,
		userRepo: userRepo,
	}
}

type CreateUserParams struct {
	Email    string      `json:"email,omitempty"    validate:"required,email"`
	Name     string      `json:"name,omitempty"     validate:"required,min=1,max=255"`
	Password string      `json:"password,omitempty" validate:"required,min=8,max=64"`
	Role     entity.Role `json:"role,omitempty"     validate:"required,min=1,max=2"`
}

func (c *CreateUser) Execute(params CreateUserParams) error {
	if err := c.val.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return validationErr
	}

	params.Email = strings.Trim(strings.ToLower(params.Email), " ")

	userWithSameEmail, err := c.userRepo.GetUserByEmail(
		context.Background(),
		params.Email,
	)
	if err != nil {
		return entity.NewErr(err)
	}
	if userWithSameEmail.ID != "" {
		return entity.ErrEmailAlreadyExists
	}

	hashedPassword, err := c.bcr.Hash(params.Password)
	if err != nil {
		return entity.NewErr(err)
	}

	repoParams := repo.CreateUserParams{}
	if err = copier.Copy(&repoParams, params); err != nil {
		return entity.NewErr(err)
	}

	repoParams.Password = hashedPassword

	if err = c.userRepo.CreateUser(context.Background(), repoParams); err != nil {
		return entity.NewErr(err)
	}

	return nil
}
