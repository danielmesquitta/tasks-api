package usecase

import (
	"testing"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/hasher"
	"github.com/danielmesquitta/tasks-api/internal/pkg/jwtutil"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestAuthenticate_Execute(t *testing.T) {
	val := validator.NewValidate()
	env := config.LoadEnv(val)
	j := jwtutil.NewJWT(env)
	bcr := hasher.NewBcrypt()

	userRepo := inmemoryrepo.NewInMemoryUserRepo()

	password := "P@ssw0rd"
	hashedPassword, err := bcr.Hash(password)
	if err != nil {
		t.Fatal(err)
	}

	user := entity.User{
		ID:        uuid.NewString(),
		Role:      entity.RoleManager,
		Name:      "John Doe",
		Email:     "johndoe@email.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userRepo.Users = append(
		userRepo.Users,
		user,
	)

	type fields struct {
		val      validator.Validator
		jwt      jwtutil.JWTManager
		bcrypt   hasher.Hasher
		userRepo *inmemoryrepo.InMemoryUserRepo
	}
	type args struct {
		params AuthenticateParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "should authenticate",
			fields: fields{
				val:      val,
				jwt:      j,
				bcrypt:   bcr,
				userRepo: userRepo,
			},
			args: args{
				params: AuthenticateParams{
					Email:    user.Email,
					Password: password,
				},
			},
			wantErr: nil,
		},
		{
			name: "should not authenticate without email",
			fields: fields{
				val:      val,
				jwt:      j,
				bcrypt:   bcr,
				userRepo: userRepo,
			},
			args: args{
				params: AuthenticateParams{
					Password: password,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not authenticate without password",
			fields: fields{
				val:      val,
				jwt:      j,
				bcrypt:   bcr,
				userRepo: userRepo,
			},
			args: args{
				params: AuthenticateParams{
					Email: user.Email,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not authenticate with invalid email",
			fields: fields{
				val:      val,
				jwt:      j,
				bcrypt:   bcr,
				userRepo: userRepo,
			},
			args: args{
				params: AuthenticateParams{
					Email:    "invalid-email",
					Password: password,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not authenticate with non-existing email",
			fields: fields{
				val:      val,
				jwt:      j,
				bcrypt:   bcr,
				userRepo: userRepo,
			},
			args: args{
				params: AuthenticateParams{
					Email:    "non-existing@email.com",
					Password: password,
				},
			},
			wantErr: entity.ErrUserEmailOrPasswordIncorrect,
		},
		{
			name: "should not authenticate with wrong password",
			fields: fields{
				val:      val,
				jwt:      j,
				bcrypt:   bcr,
				userRepo: userRepo,
			},
			args: args{
				params: AuthenticateParams{
					Email:    user.Email,
					Password: "wrong-password",
				},
			},
			wantErr: entity.ErrUserEmailOrPasswordIncorrect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := NewAuthenticate(
				tt.fields.val,
				tt.fields.jwt,
				tt.fields.bcrypt,
				tt.fields.userRepo,
			)

			gotAccessToken, gotRefreshToken, err := a.Execute(tt.args.params)
			if !testutil.IsSameErr(err, tt.wantErr) {
				t.Errorf(
					"Authenticate.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if tt.wantErr != nil {
				return
			}

			if gotAccessToken == "" {
				t.Errorf(
					"Authenticate.Execute() gotAccessToken = %v, want not empty",
					gotAccessToken,
				)
			}

			if gotRefreshToken == "" {
				t.Errorf(
					"Authenticate.Execute() gotAccessToken = %v, want not empty",
					gotAccessToken,
				)
			}
		})
	}
}
