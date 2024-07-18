package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
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

func (h TaskHandler) ListTasks(c echo.Context) error {
	tasks, err := h.listTasksUseCase.Execute(
		usecase.ListTasksParams{
			UserRole: entity.RoleManager,
			UserID:   uuid.NewString(),
		},
	)
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, tasks)
}
