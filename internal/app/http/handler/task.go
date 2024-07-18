package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/danielmesquitta/tasks-api/pkg/jwtutil"
)

type TaskHandler struct {
	listTasksUseCase *usecase.ListTasks
}

func NewTaskHandler(
	listTasksUseCase *usecase.ListTasks,
) *TaskHandler {
	return &TaskHandler{
		listTasksUseCase: listTasksUseCase,
	}
}

// @Summary List tasks
// @Description List tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {object} []entity.Task
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /tasks [get]
func (h TaskHandler) ListTasks(c echo.Context) error {
	claims, ok := c.Get("claims").(*jwtutil.UserClaims)
	if !ok {
		return entity.NewErr("invalid claims")
	}

	tasks, err := h.listTasksUseCase.Execute(
		usecase.ListTasksParams{
			UserRole: claims.Role,
			UserID:   claims.Issuer,
		},
	)
	if err != nil {
		return entity.NewErr(err)
	}
	return c.JSON(http.StatusOK, tasks)
}
