package handler

import (
	"github.com/danielmesquitta/tasks-api/internal/app/http/dto"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	createUserUseCase *usecase.CreateUser
}

func NewUserHandler(
	createUserUseCase *usecase.CreateUser,
) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
	}
}

// @Summary Create user
// @Description Create new user account (for role manager use 1 and for technician use 2)
// @Tags Users
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequestDTO true "Request body"
// @Success 201
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /users [post]
func (h *UserHandler) Create(c echo.Context) error {
	params := dto.CreateUserRequestDTO{}
	if err := c.Bind(&params); err != nil {
		return entity.NewErr(err)
	}

	useCaseParams := usecase.CreateUserParams{}
	if err := copier.Copy(&useCaseParams, params); err != nil {
		return entity.NewErr(err)
	}

	err := h.createUserUseCase.Execute(useCaseParams)
	if err != nil {
		return entity.NewErr(err)
	}

	return c.NoContent(201)
}
