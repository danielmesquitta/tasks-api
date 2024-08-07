package dto

import "github.com/danielmesquitta/tasks-api/internal/domain/entity"

type CreateUserRequestDTO struct {
	Name     string      `json:"name"`
	Role     entity.Role `json:"role"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}
