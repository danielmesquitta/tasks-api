package restapi

import (
	"go.uber.org/fx"

	"github.com/danielmesquitta/tasks-api/internal/app/restapi/handler"
	"github.com/danielmesquitta/tasks-api/internal/app/restapi/middleware"
	"github.com/danielmesquitta/tasks-api/internal/app/restapi/router"
	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/usecase"
	"github.com/danielmesquitta/tasks-api/internal/pkg/hasher"
	"github.com/danielmesquitta/tasks-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/tasks-api/internal/pkg/symcrypt"
	"github.com/danielmesquitta/tasks-api/internal/pkg/transactioner"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/broker"
	"github.com/danielmesquitta/tasks-api/internal/provider/broker/clibroker"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/mysqlrepo"
	"github.com/labstack/echo/v4"
)

func Start() {
	depsProvider := fx.Provide(
		// Config
		config.LoadEnv,

		// PKGs
		fx.Annotate(
			validator.NewValidate,
			fx.As(new(validator.Validator)),
		),
		fx.Annotate(
			symcrypt.NewAESCrypto,
			fx.As(new(symcrypt.SymmetricalEncrypter)),
		),
		fx.Annotate(
			hasher.NewBcrypt,
			fx.As(new(hasher.Hasher)),
		),
		fx.Annotate(
			jwtutil.NewJWT,
			fx.As(new(jwtutil.JWTManager)),
		),
		fx.Annotate(
			transactioner.NewSQLTransactioner,
			fx.As(new(transactioner.Transactioner)),
		),

		// Providers
		mysqlrepo.NewMySQLDBConn,
		mysqlrepo.NewMySQLQueries,
		fx.Annotate(
			mysqlrepo.NewMySQLTaskRepo,
			fx.As(new(repo.TaskRepo)),
		),
		fx.Annotate(
			mysqlrepo.NewMySQLUserRepo,
			fx.As(new(repo.UserRepo)),
		),

		fx.Annotate(
			clibroker.NewCLIMessageBroker,
			fx.As(new(broker.MessageBroker)),
		),

		// Use cases
		usecase.NewListTasks,
		usecase.NewAuthenticate,
		usecase.NewCreateUser,
		usecase.NewCreateTask,
		usecase.NewFinishTask,

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
