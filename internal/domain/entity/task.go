package entity

import "time"

type Task struct {
	ID               string     `json:"id,omitempty"`
	Summary          string     `json:"summary,omitempty"`
	AssignedToUserID string     `json:"assigned_to_user_id,omitempty"`
	CreatedByUserID  string     `json:"created_by_user_id,omitempty"`
	FinishedAt       *time.Time `json:"finished_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at,omitempty"`
}
