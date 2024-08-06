package usecase

import (
	"context"
	"testing"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/hasher"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestCreateUser_Execute(t *testing.T) {
	val := validator.NewValidate()
	bcr := hasher.NewBcrypt()

	newUserRepo := func() *inmemoryrepo.InMemoryUserRepo {
		userRepo := inmemoryrepo.NewInMemoryUserRepo()

		existingUser := entity.User{
			ID:    uuid.NewString(),
			Email: "existing-user@email.com",
		}

		userRepo.Users = append(userRepo.Users, existingUser)

		return userRepo
	}

	type fields struct {
		val      validator.Validator
		bcr      hasher.Hasher
		userRepo *inmemoryrepo.InMemoryUserRepo
	}
	type args struct {
		params CreateUserParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "should create a technician user",
			fields: fields{
				val:      val,
				bcr:      bcr,
				userRepo: newUserRepo(),
			},
			args: args{
				params: CreateUserParams{
					Name:     "John Doe",
					Email:    "johndoe@email.com",
					Password: "P@ssw0rd",
					Role:     entity.RoleTechnician,
				},
			},
			wantErr: nil,
		},
		{
			name: "should create a manager user",
			fields: fields{
				val:      val,
				bcr:      bcr,
				userRepo: newUserRepo(),
			},
			args: args{
				params: CreateUserParams{
					Name:     "John Doe",
					Email:    "johndoe@email.com",
					Password: "P@ssw0rd",
					Role:     entity.RoleManager,
				},
			},
			wantErr: nil,
		},
		{
			name: "should not create a user with invalid email",
			fields: fields{
				val:      val,
				bcr:      bcr,
				userRepo: newUserRepo(),
			},
			args: args{
				params: CreateUserParams{
					Name:     "John Doe",
					Email:    "invalid-email",
					Password: "P@ssw0rd",
					Role:     entity.RoleManager,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not create a user without a name",
			fields: fields{
				val:      val,
				bcr:      bcr,
				userRepo: newUserRepo(),
			},
			args: args{
				params: CreateUserParams{
					Email:    "johndoe@email.com",
					Password: "P@ssw0rd",
					Role:     entity.RoleManager,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not create a user without a short password",
			fields: fields{
				val:      val,
				bcr:      bcr,
				userRepo: newUserRepo(),
			},
			args: args{
				params: CreateUserParams{
					Name:     "John Doe",
					Email:    "johndoe@email.com",
					Password: "123",
					Role:     entity.RoleManager,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not create a user without a valid role",
			fields: fields{
				val:      val,
				bcr:      bcr,
				userRepo: newUserRepo(),
			},
			args: args{
				params: CreateUserParams{
					Name:     "John Doe",
					Email:    "johndoe@email.com",
					Password: "P@ssw0rd",
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not create a user with repeated email",
			fields: fields{
				val:      val,
				bcr:      bcr,
				userRepo: newUserRepo(),
			},
			args: args{
				params: CreateUserParams{
					Name:     "John Doe",
					Email:    "existing-user@email.com",
					Password: "P@ssw0rd",
					Role:     entity.RoleManager,
				},
			},
			wantErr: entity.ErrEmailAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewCreateUser(
				tt.fields.val,
				tt.fields.bcr,
				tt.fields.userRepo,
			)
			err := c.Execute(context.Background(), tt.args.params)
			if !testutil.IsSameErr(err, tt.wantErr) {
				t.Errorf(
					"CreateUser.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
		})
	}
}
