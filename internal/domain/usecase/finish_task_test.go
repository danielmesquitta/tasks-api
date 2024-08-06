package usecase

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/broker"
	"github.com/danielmesquitta/tasks-api/internal/provider/broker/clibroker"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestFinishTask_Execute(t *testing.T) {
	managerUser := entity.User{
		ID: uuid.NewString(),
	}

	technicianUser := entity.User{
		ID: uuid.NewString(),
	}

	taskRepo := inmemoryrepo.NewInMemoryTaskRepo()

	task := entity.Task{
		ID:               uuid.NewString(),
		Summary:          "Loren ipsum dolor sit amet",
		AssignedToUserID: &technicianUser.ID,
		CreatedByUserID:  managerUser.ID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	taskRepo.Tasks = append(
		taskRepo.Tasks,
		task,
	)

	type fields struct {
		validator validator.Validator
		msgBroker *clibroker.CLIMessageBroker
		taskRepo  *inmemoryrepo.InMemoryTaskRepo
	}
	type args struct {
		params FinishTaskParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "should update task finished at",
			fields: fields{
				validator: validator.NewValidate(),
				msgBroker: clibroker.NewCLIMessageBroker(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: FinishTaskParams{
					TaskID:   task.ID,
					UserID:   technicianUser.ID,
					UserRole: entity.RoleTechnician,
				},
			},
			wantErr: nil,
		},
		{
			name: "should not update task finished at if user role is not technician",
			fields: fields{
				validator: validator.NewValidate(),
				msgBroker: clibroker.NewCLIMessageBroker(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: FinishTaskParams{
					TaskID:   task.ID,
					UserID:   technicianUser.ID,
					UserRole: entity.RoleManager,
				},
			},
			wantErr: entity.ErrUserNotAllowedToFinishTask,
		},
		{
			name: "should not update task finished at if is a invalid task id",
			fields: fields{
				validator: validator.NewValidate(),
				msgBroker: clibroker.NewCLIMessageBroker(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: FinishTaskParams{
					TaskID:   "invalid-task-id",
					UserID:   technicianUser.ID,
					UserRole: entity.RoleTechnician,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not update task finished at if is a invalid user id",
			fields: fields{
				validator: validator.NewValidate(),
				msgBroker: clibroker.NewCLIMessageBroker(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: FinishTaskParams{
					TaskID:   task.ID,
					UserID:   "invalid-user-id",
					UserRole: entity.RoleTechnician,
				},
			},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not update task finished at if task does not exists",
			fields: fields{
				validator: validator.NewValidate(),
				msgBroker: clibroker.NewCLIMessageBroker(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: FinishTaskParams{
					TaskID:   uuid.NewString(),
					UserID:   technicianUser.ID,
					UserRole: entity.RoleTechnician,
				},
			},
			wantErr: entity.ErrTaskNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := NewFinishTask(
				tt.fields.validator,
				tt.fields.msgBroker,
				tt.fields.taskRepo,
			)

			sentMessages := 0
			wg := sync.WaitGroup{}
			wg.Add(1)
			if tt.wantErr == nil {
				_ = tt.fields.msgBroker.Subscribe(
					broker.TopicTaskFinished,
					func(_ []byte) {
						defer wg.Done()
						sentMessages++
					},
				)
			}

			err := f.Execute(context.Background(), tt.args.params)
			if !testutil.IsSameErr(err, tt.wantErr) {
				t.Errorf(
					"FinishTask.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}

			if tt.wantErr == nil {
				if tt.fields.taskRepo.Tasks[0].FinishedAt.IsZero() {
					t.Errorf(
						"FinishTask.Execute() taskRepo.Tasks[0].FinishedAt = %v, want not zero",
						tt.fields.taskRepo.Tasks[0].FinishedAt,
					)
				}

				wg.Wait()
				if sentMessages != 1 {
					t.Errorf(
						"FinishTask.Execute() sentMessages = %v, want 1",
						sentMessages,
					)
				}
			}
		})
	}
}
