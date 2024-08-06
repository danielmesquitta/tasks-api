package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestDeleteTask_Execute(t *testing.T) {
	newTaskRepo := func() *inmemoryrepo.InMemoryTaskRepo {
		taskRepo := inmemoryrepo.NewInMemoryTaskRepo()

		assignedToUserID := uuid.NewString()
		task := entity.Task{
			ID:               uuid.NewString(),
			Summary:          "Loren ipsum dolor sit amet",
			AssignedToUserID: &assignedToUserID,
			CreatedByUserID:  uuid.NewString(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		taskRepo.Tasks = append(
			taskRepo.Tasks,
			task,
		)

		return taskRepo
	}

	type fields struct {
		validator validator.Validator
		taskRepo  *inmemoryrepo.InMemoryTaskRepo
	}
	type args struct {
		params DeleteTaskParams
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
			task := taskRepo.Tasks[0]

			return test{
				name: "should delete a task",
				fields: fields{
					validator: validator.NewValidate(),
					taskRepo:  taskRepo,
				},
				args: args{
					params: DeleteTaskParams{
						TaskID:   task.ID,
						UserRole: entity.RoleManager,
					},
				},
				wantErr: nil,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			task := taskRepo.Tasks[0]

			return test{
				name: "should not delete a task with not allowed role",
				fields: fields{
					validator: validator.NewValidate(),
					taskRepo:  taskRepo,
				},
				args: args{
					params: DeleteTaskParams{
						TaskID:   task.ID,
						UserRole: entity.RoleTechnician,
					},
				},
				wantErr: entity.ErrUserNotAllowedToDeleteTask,
			}
		}(),
		func() test {
			taskRepo := newTaskRepo()
			task := taskRepo.Tasks[0]

			return test{
				name: "should not delete a task with invalid role",
				fields: fields{
					validator: validator.NewValidate(),
					taskRepo:  taskRepo,
				},
				args: args{
					params: DeleteTaskParams{
						TaskID:   task.ID,
						UserRole: 0,
					},
				},
				wantErr: entity.ErrValidation,
			}
		}(),
		{
			name: "should not delete a task with invalid task id",
			fields: fields{
				validator: validator.NewValidate(),
				taskRepo:  newTaskRepo(),
			},
			args: args{
				params: DeleteTaskParams{
					TaskID:   "invalid-task-id",
					UserRole: entity.RoleManager,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not delete a task with non-existing task id",
			fields: fields{
				validator: validator.NewValidate(),
				taskRepo:  newTaskRepo(),
			},
			args: args{
				params: DeleteTaskParams{
					TaskID:   uuid.NewString(),
					UserRole: entity.RoleManager,
				},
			},
			wantErr: entity.ErrTaskNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDeleteTask(
				tt.fields.validator,
				tt.fields.taskRepo,
			)
			if err := d.Execute(context.Background(), tt.args.params); !testutil.IsSameErr(
				err,
				tt.wantErr,
			) {
				t.Errorf(
					"DeleteTask.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
