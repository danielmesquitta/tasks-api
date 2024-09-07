package jwtutil

import (
	"time"

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secretKey []byte
}

func NewJWT(
	env *config.Env,
) *JWT {
	return &JWT{
		secretKey: []byte(env.JWTSecretKey),
	}
}

func (j *JWT) NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString(j.secretKey)
}

func (j *JWT) NewRefreshToken(claims jwt.RegisteredClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString(j.secretKey)
}

func (j *JWT) ValidateAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		},
	)
	if err != nil {
		return nil, entity.NewErr(err)
	}

	userClaims, ok := parsedAccessToken.Claims.(*UserClaims)
	if !ok {
		return nil, entity.NewErr("invalid claims")
	}

	if j.isExpired(&userClaims.RegisteredClaims) {
		return nil, entity.NewErr("token is expired")
	}

	return userClaims, nil
}

func (j *JWT) ValidateRefreshToken(
	refreshToken string,
) (*jwt.RegisteredClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(
		refreshToken,
		&jwt.RegisteredClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		},
	)
	if err != nil {
		return nil, entity.NewErr(err)
	}

	claims, ok := parsedRefreshToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, entity.NewErr("invalid claims")
	}

	if j.isExpired(claims) {
		return nil, entity.NewErr("token is expired")
	}

	return claims, nil
}

func (j *JWT) isExpired(claims *jwt.RegisteredClaims) bool {
	return claims.ExpiresAt.Before(time.Now())
}
