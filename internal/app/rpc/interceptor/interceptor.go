package interceptor

import (
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/jwtutil"
)

type Interceptor struct {
	AllowedRolesByMethod map[string][]entity.Role

	jwt jwtutil.JWTManager
}

func NewInterceptor(
	jwt jwtutil.JWTManager,
) *Interceptor {
	return &Interceptor{
		jwt: jwt,
	}
}
