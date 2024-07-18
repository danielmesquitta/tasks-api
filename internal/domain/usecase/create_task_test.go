package usecase

import (
	"strings"
	"testing"

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/pkg/crypto"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestCreateTask_Execute(t *testing.T) {
	val := validator.NewValidator()
	env := config.LoadEnv(val)

	userRepo := inmemoryrepo.NewInMemoryUserRepo()

	managerUser := entity.User{
		ID:   uuid.NewString(),
		Role: entity.RoleManager,
	}

	technicianUser := entity.User{
		ID:   uuid.NewString(),
		Role: entity.RoleTechnician,
	}

	userRepo.Users = append(
		userRepo.Users,
		managerUser,
		technicianUser,
	)

	type fields struct {
		validator *validator.Validator
		crypto    *crypto.Crypto
		taskRepo  *inmemoryrepo.InMemoryTaskRepo
		userRepo  *inmemoryrepo.InMemoryUserRepo
	}
	type args struct {
		params CreateTaskParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantID  string
		wantErr error
	}{
		{
			name: "should create a task",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo:  userRepo,
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  managerUser.ID,
					AssignedToUserID: technicianUser.ID,
				},
			},
			wantErr: nil,
		},
		{
			name: "should not create a task if user role is not manager",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo:  userRepo,
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleTechnician,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  technicianUser.ID,
					AssignedToUserID: technicianUser.ID,
				},
			},
			wantErr: entity.ErrUserNotAllowedToCreateTask,
		},
		{
			name: "should not create a task if user is trying to pass as a manager",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo:  userRepo,
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  technicianUser.ID,
					AssignedToUserID: technicianUser.ID,
				},
			},
			wantErr: entity.ErrUserNotAllowedToCreateTask,
		},
		{
			name: "should not create a task if created by user is not found",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo: func() *inmemoryrepo.InMemoryUserRepo {
					customUserRepo := inmemoryrepo.NewInMemoryUserRepo()
					customUserRepo.Users = append(
						customUserRepo.Users,
						technicianUser,
					)
					return customUserRepo
				}(),
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  managerUser.ID,
					AssignedToUserID: technicianUser.ID,
				},
			},
			wantErr: entity.ErrCreatedByUserNotFound,
		},
		{
			name: "should not create a task if assigned to user is not found",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo: func() *inmemoryrepo.InMemoryUserRepo {
					customUserRepo := inmemoryrepo.NewInMemoryUserRepo()
					customUserRepo.Users = append(
						customUserRepo.Users,
						managerUser,
					)
					return customUserRepo
				}(),
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  managerUser.ID,
					AssignedToUserID: technicianUser.ID,
				},
			},
			wantErr: entity.ErrAssignToUserNotFound,
		},
		{
			name: "should not create a task if assigned to user is not a technician",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo:  userRepo,
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  managerUser.ID,
					AssignedToUserID: managerUser.ID,
				},
			},
			wantErr: entity.ErrInvalidRoleForAssignedUser,
		},
		{
			name: "should not create a task if summary is greater than 2500 characters",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo:  userRepo,
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          strings.Repeat("a", 2501),
					CreatedByUserID:  managerUser.ID,
					AssignedToUserID: technicianUser.ID,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not create a task if created by user id is invalid",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo:  userRepo,
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  "invalid",
					AssignedToUserID: technicianUser.ID,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not create a task if assigned to user id is invalid",
			fields: fields{
				validator: val,
				crypto:    crypto.NewCrypto(env),
				taskRepo:  inmemoryrepo.NewInMemoryTaskRepo(),
				userRepo:  userRepo,
			},
			args: args{
				params: CreateTaskParams{
					UserRole:         entity.RoleManager,
					Summary:          "Loren ipsum dolor sit amet",
					CreatedByUserID:  managerUser.ID,
					AssignedToUserID: "invalid",
				},
			},
			wantErr: entity.ErrValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCreateTask(
				tt.fields.validator,
				tt.fields.crypto,
				tt.fields.taskRepo,
				tt.fields.userRepo,
			)
			gotID, err := c.Execute(tt.args.params)
			if !testutil.IsSameErr(err, tt.wantErr) {
				t.Errorf(
					"CreateTask.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if tt.wantErr != nil {
				return
			}

			tt.wantID = tt.fields.taskRepo.Tasks[0].ID

			if gotID != tt.wantID {
				t.Errorf(
					"CreateTask.Execute() = %v, want %v",
					gotID,
					tt.wantID,
				)
			}

		})
	}
}
