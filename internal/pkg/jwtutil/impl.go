package jwtutil

import (
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

func (j *JWT) ParseAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
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

	return userClaims, nil
}

func (j *JWT) ParseRefreshToken(
	refreshToken string,
) (*jwt.RegisteredClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(
		refreshToken,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
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

	return claims, nil
}
