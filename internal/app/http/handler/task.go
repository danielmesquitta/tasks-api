package handler

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"

	"github.com/danielmesquitta/tasks-api/internal/app/http/dto"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/danielmesquitta/tasks-api/internal/pkg/jwtutil"
)

type TaskHandler struct {
	createTaskUseCase *usecase.CreateTask
	finishTaskUseCase *usecase.FinishTask
	listTasksUseCase  *usecase.ListTasks
}

func NewTaskHandler(
	createTaskUseCase *usecase.CreateTask,
	finishTaskUseCase *usecase.FinishTask,
	listTasksUseCase *usecase.ListTasks,
) *TaskHandler {
	return &TaskHandler{
		createTaskUseCase: createTaskUseCase,
		finishTaskUseCase: finishTaskUseCase,
		listTasksUseCase:  listTasksUseCase,
	}
}

// @Summary Create task
// @Description Create new task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body dto.CreateTaskRequestDTO true "Request body"
// @Success 201
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 404 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /tasks [post]
func (h *TaskHandler) Create(c echo.Context) error {
	claims, ok := c.Get("claims").(*jwtutil.UserClaims)
	if !ok {
		return entity.NewErr("invalid claims")
	}

	params := dto.CreateTaskRequestDTO{}
	if err := c.Bind(&params); err != nil {
		return entity.NewErr(err)
	}

	useCaseParams := usecase.CreateTaskParams{}
	if err := copier.Copy(&useCaseParams, params); err != nil {
		return entity.NewErr(err)
	}

	useCaseParams.UserRole = claims.Role
	useCaseParams.CreatedByUserID = claims.Issuer

	err := h.createTaskUseCase.Execute(useCaseParams)
	if err != nil {
		return entity.NewErr(err)
	}

	return c.NoContent(http.StatusCreated)
}

// @Summary Finish task
// @Description Mark task as finished
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 404 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /tasks/{id}/finished [patch]
func (h *TaskHandler) Finish(c echo.Context) error {
	claims, ok := c.Get("claims").(*jwtutil.UserClaims)
	if !ok {
		return entity.NewErr("invalid claims")
	}

	useCaseParams := usecase.FinishTaskParams{
		TaskID:   c.Param("id"),
		UserID:   claims.Issuer,
		UserRole: claims.Role,
	}

	err := h.finishTaskUseCase.Execute(useCaseParams)
	if err != nil {
		return entity.NewErr(err)
	}

	return c.NoContent(http.StatusOK)
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
func (h TaskHandler) List(c echo.Context) error {
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
