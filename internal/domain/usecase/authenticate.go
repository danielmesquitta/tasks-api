package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/pkg/cryptoutil"
	"github.com/danielmesquitta/tasks-api/pkg/jwtutil"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticate struct {
	val      *validator.Validator
	jwt      *jwtutil.JWT
	bcrypt   *cryptoutil.Bcrypt
	userRepo repo.UserRepo
}

func NewAuthenticate(
	val *validator.Validator,
	jwt *jwtutil.JWT,
	bcrypt *cryptoutil.Bcrypt,
	userRepo repo.UserRepo,
) *Authenticate {
	return &Authenticate{
		val:      val,
		jwt:      jwt,
		bcrypt:   bcrypt,
		userRepo: userRepo,
	}
}

type AuthenticateParams struct {
	Email    string `json:"email,omitempty"    validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (a *Authenticate) Execute(
	params AuthenticateParams,
) (accessToken, refreshToken string, err error) {
	if err = a.val.Validate(params); err != nil {
		validationErr := entity.ErrValidation
		validationErr.Message = err.Error()
		return "", "", validationErr
	}

	user, err := a.userRepo.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		return "", "", entity.NewErr(err)
	}

	if user.ID == "" {
		return "", "", entity.ErrUserEmailOrPasswordIncorrect
	}

	if !a.bcrypt.Match(params.Password, user.Password) {
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
