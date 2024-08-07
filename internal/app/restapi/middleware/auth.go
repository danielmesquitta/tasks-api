package middleware

import (
	"net/http"
	"strings"

	"github.com/danielmesquitta/tasks-api/internal/app/restapi/dto"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) EnsureAuthenticated(
	next echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(
				http.StatusUnauthorized,
				dto.ErrorResponseDTO{Message: "missing or malformed token"},
			)
		}

		// Split the header to get the token part
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.JSON(
				http.StatusUnauthorized,
				dto.ErrorResponseDTO{Message: "invalid token format"},
			)
		}

		accessToken := parts[1]

		// Parse and validate the token
		claims, err := m.jwt.ValidateAccessToken(accessToken)
		if err != nil {
			return c.JSON(
				http.StatusUnauthorized,
				dto.ErrorResponseDTO{Message: "invalid token"},
			)
		}

		// Set the claims in the context
		c.Set("claims", claims)

		// Token is valid, proceed with the request
		return next(c)
	}
}
