package middleware

import (
	"github.com/danielmesquitta/tasks-api/pkg/jwtutil"
	"github.com/danielmesquitta/tasks-api/pkg/logger"
)

type Middleware struct {
	log *logger.Logger
	jwt *jwtutil.JWT
}

func NewMiddleware(
	log *logger.Logger,
	jwt *jwtutil.JWT,
) *Middleware {
	return &Middleware{
		log: log,
		jwt: jwt,
	}
}
