package router

import (
	_ "github.com/danielmesquitta/tasks-api/docs"
	"github.com/danielmesquitta/tasks-api/internal/app/http/handler"
	"github.com/danielmesquitta/tasks-api/internal/config"

	"github.com/labstack/echo/v4"
)

type Router struct {
	env         *config.Env
	taskHandler *handler.TaskHandler
}

func NewRouter(
	env *config.Env,
	taskHandler *handler.TaskHandler,
) *Router {
	return &Router{
		env:         env,
		taskHandler: taskHandler,
	}
}

func (r *Router) Register(
	app *echo.Echo,
) {
	basePath := "/api/v1"
	apiV1 := app.Group(basePath)

	apiV1.GET("/tasks", r.taskHandler.ListTasks)
}
