package middleware

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (m *Middleware) BasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	middlewareFunc := middleware.BasicAuth(
		func(username, password string, _ echo.Context) (bool, error) {
			usernameMatches := subtle.ConstantTimeCompare(
				[]byte(username),
				[]byte(m.env.BasicAuthUsername),
			) == 1
			passwordMatches := subtle.ConstantTimeCompare(
				[]byte(password),
				[]byte(m.env.BasicAuthPassword),
			) == 1

			if usernameMatches && passwordMatches {
				return true, nil
			}

			return false, nil
		},
	)

	return middlewareFunc(next)
}
