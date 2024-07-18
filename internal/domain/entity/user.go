package entity

import "time"

type Role byte

const (
	RoleManager Role = iota + 1
	RoleTechnician
)

type User struct {
	ID        string    `json:"id"`
	Role      Role      `json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
