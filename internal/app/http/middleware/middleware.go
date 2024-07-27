package middleware

import (
	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/pkg/jwtutil"
)

type Middleware struct {
	env *config.Env
	jwt jwtutil.JWTManager
}

func NewMiddleware(
	env *config.Env,
	jwt jwtutil.JWTManager,
) *Middleware {
	return &Middleware{
		env: env,
		jwt: jwt,
	}
}
