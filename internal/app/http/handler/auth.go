package handler

import (
	"net/http"

	"github.com/danielmesquitta/tasks-api/internal/app/http/dto"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authenticateUseCase *usecase.Authenticate
}

func NewAuthHandler(
	authenticateUseCase *usecase.Authenticate,
) *AuthHandler {
	return &AuthHandler{
		authenticateUseCase: authenticateUseCase,
	}
}

// @Summary Login
// @Description Authenticate user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AuthenticateRequestDTO true "Request body"
// @Success 200 {object} dto.AuthenticateResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /auth/login [post]
func (a *AuthHandler) Login(c echo.Context) error {
	requestData := dto.AuthenticateRequestDTO{}
	if err := c.Bind(&requestData); err != nil {
		return entity.NewErr(err)
	}

	useCaseParams := usecase.AuthenticateParams{}
	if err := copier.Copy(&useCaseParams, requestData); err != nil {
		return entity.NewErr(err)
	}

	accessToken, refreshToken, err := a.authenticateUseCase.Execute(
		c.Request().Context(),
		useCaseParams,
	)
	if err != nil {
		return entity.NewErr(err)
	}

	return c.JSON(http.StatusCreated, dto.AuthenticateResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
