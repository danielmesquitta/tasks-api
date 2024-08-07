package dto

type CreateTaskRequestDTO struct {
	Summary          string `json:"summary,omitempty"`
	AssignedToUserID string `json:"assigned_to_user_id,omitempty"`
}
