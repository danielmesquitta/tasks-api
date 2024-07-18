package http

import (
	"github.com/danielmesquitta/tasks-api/internal/app/http/handler"
	"github.com/danielmesquitta/tasks-api/internal/app/http/router"
	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/mysqlrepo"
	"github.com/danielmesquitta/tasks-api/pkg/crypto"
	"github.com/danielmesquitta/tasks-api/pkg/logger"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func Start() {
	depsProvider := fx.Provide(
		// Config
		config.LoadEnv,

		// PKGs
		validator.NewValidator,
		logger.NewLogger,
		crypto.NewCrypto,

		// Providers
		mysqlrepo.NewMySQLDBConn,
		fx.Annotate(
			mysqlrepo.NewMySQLTaskRepo,
			fx.As(new(repo.TaskRepo)),
		),

		// Use cases
		usecase.NewListTasks,

		// Handlers
		handler.NewTaskHandler,

		// Router
		router.NewRouter,

		// App
		NewApp,
	)

	container := fx.New(
		depsProvider,
		fx.Invoke(func(*echo.Echo) {}),
	)

	container.Run()
}
