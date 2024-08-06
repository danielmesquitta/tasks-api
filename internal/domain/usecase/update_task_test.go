package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/symcrypt"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestUpdateTask_Execute(t *testing.T) {
	val := validator.NewValidate()
	env := config.LoadEnv(val)
	symCrypto := symcrypt.NewAESCrypto(env)

	technicianUser := entity.User{
		ID:   uuid.NewString(),
		Role: entity.RoleTechnician,
	}

	managerUser := entity.User{
		ID:   uuid.NewString(),
		Role: entity.RoleManager,
	}

	newUserRepo := func() *inmemoryrepo.InMemoryUserRepo {
		userRepo := inmemoryrepo.NewInMemoryUserRepo()
		userRepo.Users = append(
			userRepo.Users,
			technicianUser,
			managerUser,
		)
		return userRepo
	}

	beforeUpdateSummary := "Loren ipsum dolor sit amet"
	newTaskRepo := func() *inmemoryrepo.InMemoryTaskRepo {
		taskRepo := inmemoryrepo.NewInMemoryTaskRepo()
		task := entity.Task{
			ID:              uuid.NewString(),
			Summary:         beforeUpdateSummary,
			CreatedByUserID: uuid.NewString(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		taskRepo.Tasks = append(
			taskRepo.Tasks,
			task,
		)
		return taskRepo
	}

	type fields struct {
		validator validator.Validator
		symCrypto symcrypt.SymmetricalEncrypter
		taskRepo  *inmemoryrepo.InMemoryTaskRepo
		userRepo  *inmemoryrepo.InMemoryUserRepo
	}
	type args struct {
		params UpdateTaskParams
	}
	type test struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}
	tests := []test{
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should update a task",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               taskRepo.Tasks[0].ID,
						UserID:           managerUser.ID,
						UserRole:         entity.RoleManager,
						Summary:          "Loren ipsum",
						AssignedToUserID: &technicianUser.ID,
					},
				},
				wantErr: nil,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should update a task without changing the assigned user",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:       taskRepo.Tasks[0].ID,
						UserID:   managerUser.ID,
						UserRole: entity.RoleManager,
						Summary:  "Loren ipsum",
					},
				},
				wantErr: nil,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should not update a task with invalid id",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               "invalid-id",
						UserID:           managerUser.ID,
						UserRole:         entity.RoleManager,
						Summary:          "Loren ipsum",
						AssignedToUserID: &technicianUser.ID,
					},
				},
				wantErr: entity.ErrValidation,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should not update a task with invalid user id",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               taskRepo.Tasks[0].ID,
						UserID:           "invalid-user-id",
						UserRole:         entity.RoleManager,
						Summary:          "Loren ipsum",
						AssignedToUserID: &technicianUser.ID,
					},
				},
				wantErr: entity.ErrValidation,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should not update a task with invalid user role",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               taskRepo.Tasks[0].ID,
						UserID:           managerUser.ID,
						UserRole:         0,
						Summary:          "Loren ipsum",
						AssignedToUserID: &technicianUser.ID,
					},
				},
				wantErr: entity.ErrValidation,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should not update a task with summary greater than 2500 characters",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               taskRepo.Tasks[0].ID,
						UserID:           managerUser.ID,
						UserRole:         entity.RoleManager,
						Summary:          strings.Repeat("a", 2501),
						AssignedToUserID: &technicianUser.ID,
					},
				},
				wantErr: entity.ErrValidation,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			invalidUserID := "invalid-user-id"
			return test{
				name: "should not update a task with invalid assigned user id",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               taskRepo.Tasks[0].ID,
						UserID:           managerUser.ID,
						UserRole:         entity.RoleManager,
						Summary:          "Lorem ipsum",
						AssignedToUserID: &invalidUserID,
					},
				},
				wantErr: entity.ErrValidation,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should not update non-existent task",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               uuid.NewString(),
						UserID:           managerUser.ID,
						UserRole:         entity.RoleManager,
						Summary:          "Loren ipsum",
						AssignedToUserID: &technicianUser.ID,
					},
				},
				wantErr: entity.ErrTaskNotFound,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			nonExistingUserID := uuid.NewString()
			return test{
				name: "should not update task with non-existent assigned user",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               taskRepo.Tasks[0].ID,
						UserID:           managerUser.ID,
						UserRole:         entity.RoleManager,
						Summary:          "Loren ipsum",
						AssignedToUserID: &nonExistingUserID,
					},
				},
				wantErr: entity.ErrUserNotFound,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should not update task assigned user if user role is not manager",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:               taskRepo.Tasks[0].ID,
						UserID:           technicianUser.ID,
						UserRole:         entity.RoleTechnician,
						Summary:          "Loren ipsum",
						AssignedToUserID: &technicianUser.ID,
					},
				},
				wantErr: entity.ErrUserNotAllowedToUpdateAssignedUser,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			userRepo := newUserRepo()
			return test{
				name: "should not update task assigned to other technician",
				fields: fields{
					validator: val,
					symCrypto: symCrypto,
					taskRepo:  taskRepo,
					userRepo:  userRepo,
				},
				args: args{
					params: UpdateTaskParams{
						ID:       taskRepo.Tasks[0].ID,
						UserID:   uuid.NewString(),
						UserRole: entity.RoleTechnician,
						Summary:  "Loren ipsum",
					},
				},
				wantErr: entity.ErrUserNotAllowedToUpdateTask,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			u := NewUpdateTask(
				tt.fields.validator,
				tt.fields.symCrypto,
				tt.fields.taskRepo,
				tt.fields.userRepo,
			)
			err := u.Execute(context.Background(), tt.args.params)
			if !testutil.IsSameErr(
				err,
				tt.wantErr,
			) {
				t.Errorf(
					"UpdateTask.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}

			if err != nil {
				return
			}

			task := tt.fields.taskRepo.Tasks[0]
			if !testutil.CompareAsPtr(
				task.AssignedToUserID,
				tt.args.params.AssignedToUserID,
			) {
				t.Errorf(
					"UpdateTask.Execute() task.AssignedToUserID = %v, want %v",
					*task.AssignedToUserID,
					*tt.args.params.AssignedToUserID,
				)
			}

			if task.Summary == beforeUpdateSummary &&
				tt.args.params.Summary != beforeUpdateSummary {
				t.Errorf(
					"UpdateTask.Execute() task.Summary = %v, want not %v",
					task.Summary,
					tt.args.params.Summary,
				)
			}
		})
	}
}
