package http

import (
	"github.com/danielmesquitta/tasks-api/internal/app/http/handler"
	"github.com/danielmesquitta/tasks-api/internal/app/http/middleware"
	"github.com/danielmesquitta/tasks-api/internal/app/http/router"
	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/mysqlrepo"
	"github.com/danielmesquitta/tasks-api/pkg/cryptoutil"
	"github.com/danielmesquitta/tasks-api/pkg/jwtutil"
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
		cryptoutil.NewAESCrypto,
		cryptoutil.NewBcrypt,
		jwtutil.NewJWT,

		// Providers
		mysqlrepo.NewMySQLDBConn,
		fx.Annotate(
			mysqlrepo.NewMySQLTaskRepo,
			fx.As(new(repo.TaskRepo)),
		),
		fx.Annotate(
			mysqlrepo.NewMySQLUserRepo,
			fx.As(new(repo.UserRepo)),
		),

		// Use cases
		usecase.NewListTasks,
		usecase.NewAuthenticate,
		usecase.NewCreateUser,

		// Handlers
		handler.NewAuthHandler,
		handler.NewUserHandler,
		handler.NewTaskHandler,

		// Middleware
		middleware.NewMiddleware,

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
