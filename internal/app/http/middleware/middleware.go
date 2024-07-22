package middleware

import (
	"github.com/danielmesquitta/tasks-api/internal/pkg/jwtutil"
)

type Middleware struct {
	jwt jwtutil.JWTManager
}

func NewMiddleware(
	jwt jwtutil.JWTManager,
) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}
