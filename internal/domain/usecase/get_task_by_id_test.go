package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/danielmesquitta/tasks-api/internal/pkg/symcrypt"
	"github.com/danielmesquitta/tasks-api/internal/pkg/validator"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo"
	"github.com/danielmesquitta/tasks-api/internal/provider/repo/inmemoryrepo"
	"github.com/danielmesquitta/tasks-api/test/testutil"
	"github.com/google/uuid"
)

func TestGetTaskByID_Execute(t *testing.T) {
	val := validator.NewValidate()
	symCrypto := symcrypt.NewAESCrypto(config.LoadEnv(val))

	managerUserID := uuid.NewString()
	technicianUserID := uuid.NewString()

	taskRepo := inmemoryrepo.NewInMemoryTaskRepo()

	summary := "Loren Ipsum"
	encryptedSummary, err := symCrypto.Encrypt(summary)
	if err != nil {
		t.Errorf("Error encrypting summary: %v", err)
	}
	task := entity.Task{
		ID:               uuid.NewString(),
		Summary:          encryptedSummary,
		AssignedToUserID: &technicianUserID,
		CreatedByUserID:  managerUserID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	taskWithoutAssignedUser := entity.Task{
		ID:              uuid.NewString(),
		Summary:         encryptedSummary,
		CreatedByUserID: managerUserID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	getDecryptedTask := func() entity.Task {
		t := task
		t.Summary = summary
		return t
	}

	taskRepo.Tasks = append(taskRepo.Tasks, task, taskWithoutAssignedUser)

	type fields struct {
		validator validator.Validator
		symCrypto symcrypt.SymmetricalEncrypter
		taskRepo  repo.TaskRepo
	}
	type args struct {
		params GetTaskByIDParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Task
		wantErr error
	}{
		{
			name: "should return any task if the user is a manager",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       task.ID,
					UserID:   managerUserID,
					UserRole: entity.RoleManager,
				},
			},
			want:    getDecryptedTask(),
			wantErr: nil,
		},
		{
			name: "should return assigned task if user is a technician",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       task.ID,
					UserID:   technicianUserID,
					UserRole: entity.RoleTechnician,
				},
			},
			want:    getDecryptedTask(),
			wantErr: nil,
		},
		{
			name: "should not return task if id is invalid",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       "invalid-id",
					UserID:   technicianUserID,
					UserRole: entity.RoleTechnician,
				},
			},
			want:    entity.Task{},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not return task if user id is invalid",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       task.ID,
					UserID:   "invalid-user-id",
					UserRole: entity.RoleTechnician,
				},
			},
			want:    entity.Task{},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not return task if user role is invalid",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       task.ID,
					UserID:   managerUserID,
					UserRole: 0,
				},
			},
			want:    entity.Task{},
			wantErr: entity.ErrValidation,
		},
		{
			name: "should not return non-existing task",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       uuid.NewString(),
					UserID:   managerUserID,
					UserRole: entity.RoleManager,
				},
			},
			want:    entity.Task{},
			wantErr: entity.ErrTaskNotFound,
		},
		{
			name: "should not return task if user is technician and task is not assigned to anyone",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       taskWithoutAssignedUser.ID,
					UserID:   technicianUserID,
					UserRole: entity.RoleTechnician,
				},
			},
			want:    entity.Task{},
			wantErr: entity.ErrUserNotAllowedToViewTask,
		},
		{
			name: "should not return task if user is technician and task is not assigned to him",
			fields: fields{
				validator: val,
				symCrypto: symCrypto,
				taskRepo:  taskRepo,
			},
			args: args{
				params: GetTaskByIDParams{
					ID:       task.ID,
					UserID:   uuid.NewString(),
					UserRole: entity.RoleTechnician,
				},
			},
			want:    entity.Task{},
			wantErr: entity.ErrUserNotAllowedToViewTask,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := NewGetTaskByID(
				tt.fields.validator,
				tt.fields.symCrypto,
				tt.fields.taskRepo,
			)
			got, err := u.Execute(context.Background(), tt.args.params)
			if !testutil.IsSameErr(err, tt.wantErr) {
				t.Errorf(
					"GetTaskByID.Execute() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskByID.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
