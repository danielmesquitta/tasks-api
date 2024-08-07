package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/danielmesquitta/tasks-api/docs" // swagger docs

	"github.com/danielmesquitta/tasks-api/internal/app/restapi/handler"
	mid "github.com/danielmesquitta/tasks-api/internal/app/restapi/middleware"
	"github.com/danielmesquitta/tasks-api/internal/config"
)

type Router struct {
	env         *config.Env
	mid         *mid.Middleware
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
	taskHandler *handler.TaskHandler
}

func NewRouter(
	env *config.Env,
	mid *mid.Middleware,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	taskHandler *handler.TaskHandler,
) *Router {
	return &Router{
		env:         env,
		mid:         mid,
		authHandler: authHandler,
		userHandler: userHandler,
		taskHandler: taskHandler,
	}
}

func (r *Router) Register(
	app *echo.Echo,
) {
	basePath := "/api/v1"
	apiV1 := app.Group(basePath)

	apiV1.GET("/docs/*", echoSwagger.WrapHandler)

	apiV1.POST("/users", r.userHandler.Create, r.mid.BasicAuth)

	apiV1.POST("/auth/login", r.authHandler.Login)

	apiV1.POST("/tasks", r.taskHandler.Create, r.mid.EnsureAuthenticated)
	apiV1.PATCH(
		"/tasks/:id/finished",
		r.taskHandler.Finish,
		r.mid.EnsureAuthenticated,
	)
	apiV1.GET("/tasks", r.taskHandler.List, r.mid.EnsureAuthenticated)
}
