package jwtutil

import (
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	NewAccessToken(claims UserClaims) (accessToken string, err error)
	NewRefreshToken(
		claims jwt.RegisteredClaims,
	) (refreshToken string, err error)
	ParseAccessToken(accessToken string) (*UserClaims, error)
	ParseRefreshToken(refreshToken string) (*jwt.RegisteredClaims, error)
}

type UserClaims struct {
	Role entity.Role `json:"role,omitempty"`
	jwt.RegisteredClaims
}
