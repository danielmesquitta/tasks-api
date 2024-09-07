package entity

import (
	"fmt"
	"runtime/debug"
)

type Err struct {
	Message    string  `json:"message,omitempty"`
	StackTrace string  `json:"stack_trace,omitempty"`
	Type       ErrType `json:"type,omitempty"`
}

type ErrType string

const (
	ErrTypeUnknown      ErrType = "unknown"
	ErrTypeNotFound     ErrType = "not_found"
	ErrTypeUnauthorized ErrType = "unauthorized"
	ErrTypeForbidden    ErrType = "forbidden"
	ErrTypeValidation   ErrType = "validation_error"
)

func newErr(err any, errType ErrType) *Err {
	switch v := err.(type) {
	case *Err:
		return v
	case error:
		return &Err{
			Message:    v.Error(),
			StackTrace: string(debug.Stack()),
			Type:       errType,
		}
	case string:
		return &Err{
			Message:    v,
			StackTrace: string(debug.Stack()),
			Type:       errType,
		}
	default:
		panic("trying to create an Err with an unsupported type")
	}
}

// NewErr creates a new Err instance from either an error or a string,
// and sets the Type flag to unknown. This is useful when you want to
// create an error that is not expected to happen, and you want to
// log it with stack tracing.
func NewErr(err any) *Err {
	return newErr(err, ErrTypeUnknown)
}

func (e *Err) Error() string {
	return e.Message
}

func (e *Err) ErrorWithStackTrace() string {
	return fmt.Sprintf("%s\n\n%s", e.Message, e.StackTrace)
}

var (
	ErrUserNotFound = newErr(
		"user not found",
		ErrTypeNotFound,
	)
	ErrCreatedByUserNotFound = newErr(
		"user creating task was not found",
		ErrTypeNotFound,
	)
	ErrAssignToUserNotFound = newErr(
		"user assigned to task was not found",
		ErrTypeNotFound,
	)
	ErrTaskNotFound = newErr(
		"task not found",
		ErrTypeNotFound,
	)
	ErrValidation = newErr(
		"validation error",
		ErrTypeValidation,
	)
	ErrUserNotAllowedToCreateTask = newErr(
		"only users with the role manager can create tasks",
		ErrTypeForbidden,
	)
	ErrInvalidRoleForAssignedUser = newErr(
		"only users with the role technician can be assigned to tasks",
		ErrTypeForbidden,
	)
	ErrUserNotAllowedToFinishTask = newErr(
		"only users with the role technician can finish tasks",
		ErrTypeForbidden,
	)
	ErrUserEmailOrPasswordIncorrect = newErr(
		"email or password is incorrect",
		ErrTypeUnauthorized,
	)
	ErrEmailAlreadyExists = newErr(
		"email is already registered",
		ErrTypeValidation,
	)
	ErrUserNotAllowedToDeleteTask = newErr(
		"only users with the role manager can delete tasks",
		ErrTypeForbidden,
	)
	ErrUserNotAllowedToUpdateTask = newErr(
		"only users with the role of manager or those assigned to this task can update it",
		ErrTypeForbidden,
	)
	ErrUserNotAllowedToUpdateAssignedUser = newErr(
		"only managers can update the assigned user of a task",
		ErrTypeForbidden,
	)
	ErrUserNotAllowedToViewTask = newErr(
		"only users with the role of manager or those assigned to this task can view it",
		ErrTypeForbidden,
	)
)

var _ error = (*Err)(nil)
