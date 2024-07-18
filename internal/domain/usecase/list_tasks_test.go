package usecase

import (
	"reflect"
	"testing"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/pkg/crypto"
	"github.com/danielmesquitta/tasks-api/pkg/validator"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestListTasks_Execute(t *testing.T) {
	taskRepo := inmemoryrepo.NewInMemoryTaskRepo()

	managerID := uuid.NewString()
	firstTechnicianID := uuid.NewString()
	secondTechnicianID := uuid.NewString()

	task1 := entity.Task{
		ID:               uuid.NewString(),
		Summary:          "Lorem ipsum dolor sit amet",
		AssignedToUserID: firstTechnicianID,
		CreatedByUserID:  managerID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	task2 := entity.Task{
		ID:               uuid.NewString(),
		Summary:          "Lorem ipsum dolor sit amet",
		AssignedToUserID: firstTechnicianID,
		CreatedByUserID:  managerID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	task3 := entity.Task{
		ID:               uuid.NewString(),
		Summary:          "Lorem ipsum dolor sit amet",
		AssignedToUserID: secondTechnicianID,
		CreatedByUserID:  managerID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	taskRepo.Tasks = append(
		taskRepo.Tasks,
		task1,
		task2,
		task3,
	)

	type fields struct {
		validator *validator.Validator
		crypto    *crypto.Crypto
		taskRepo  repo.TaskRepo
	}
	type args struct {
		params ListTasksParams
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantTasks []entity.Task
		wantErr   error
	}{
		{
			name: "should list all tasks for the manager",
			fields: fields{
				validator: validator.NewValidator(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: ListTasksParams{
					UserRole: entity.RoleManager,
					UserID:   managerID,
				},
			},
			wantTasks: []entity.Task{
				task1,
				task2,
				task3,
			},
			wantErr: nil,
		},
		{
			name: "should list only the tasks assigned to the first technician",
			fields: fields{
				validator: validator.NewValidator(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: ListTasksParams{
					UserRole: entity.RoleTechnician,
					UserID:   firstTechnicianID,
				},
			},
			wantTasks: []entity.Task{
				task1,
				task2,
			},
			wantErr: nil,
		},
		{
			name: "should list only the tasks assigned to the second technician",
			fields: fields{
				validator: validator.NewValidator(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: ListTasksParams{
					UserRole: entity.RoleTechnician,
					UserID:   secondTechnicianID,
				},
			},
			wantTasks: []entity.Task{
				task3,
			},
			wantErr: nil,
		},
		{
			name: "should not list task if invalid id is provided",
			fields: fields{
				validator: validator.NewValidator(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: ListTasksParams{
					UserRole: entity.RoleTechnician,
					UserID:   "invalid-id",
				},
			},
			wantTasks: nil,
			wantErr:   entity.ErrValidation,
		},
		{
			name: "should not list task if invalid role is provided",
			fields: fields{
				validator: validator.NewValidator(),
				taskRepo:  taskRepo,
			},
			args: args{
				params: ListTasksParams{
					UserRole: 0,
					UserID:   managerID,
				},
			},
			wantTasks: nil,
			wantErr:   entity.ErrValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListTasks(
				tt.fields.validator,
				tt.fields.crypto,
				tt.fields.taskRepo,
			)

			gotTasks, err := l.Execute(tt.args.params)
			if !testutil.IsSameErr(err, tt.wantErr) {
				t.Errorf(
					"ListTasks.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(gotTasks, tt.wantTasks) {
				t.Errorf(
					"ListTasks.Execute() = %v, want %v",
					gotTasks,
					tt.wantTasks,
				)
			}
		})
	}
}
