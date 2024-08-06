package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/hasher"
	"github.com/danielmesquitta/tasks-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticate struct {
	val      validator.Validator
	jwt      jwtutil.JWTManager
	hasher   hasher.Hasher
	userRepo repo.UserRepo
}

func NewAuthenticate(
	val validator.Validator,
	jwt jwtutil.JWTManager,
	hasher hasher.Hasher,
	userRepo repo.UserRepo,
) *Authenticate {
	return &Authenticate{
		val:      val,
		jwt:      jwt,
		hasher:   hasher,
		userRepo: userRepo,
	}
}

type AuthenticateParams struct {
	Email    string `json:"email,omitempty"    validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (a *Authenticate) Execute(
	ctx context.Context,
	params AuthenticateParams,
) (accessToken, refreshToken string, err error) {
	if err = a.val.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return "", "", validationErr
	}

	user, err := a.userRepo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return "", "", entity.NewErr(err)
	}

	if user.ID == "" {
		return "", "", entity.ErrUserEmailOrPasswordIncorrect
	}

	if !a.hasher.Match(params.Password, user.Password) {
		return "", "", entity.ErrUserEmailOrPasswordIncorrect
	}

	accessToken, err = a.jwt.NewAccessToken(jwtutil.UserClaims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	if err != nil {
		return "", "", entity.NewErr(err)
	}

	refreshToken, err = a.jwt.NewRefreshToken(jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	})
	if err != nil {
		return "", "", entity.NewErr(err)
	}

	return accessToken, refreshToken, nil
}
