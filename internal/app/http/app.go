package http

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/danielmesquitta/tasks-api/internal/app/http/router"
	"github.com/danielmesquitta/tasks-api/internal/config"
)

func NewApp(
	lc fx.Lifecycle,
	env *config.Env,
	router *router.Router,
) *echo.Echo {
	app := echo.New()

	router.Register(app)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := app.Start(":" + env.Port); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown(context.Background())
		},
	})

	return app
}