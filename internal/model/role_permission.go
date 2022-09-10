package model

import "time"

type RolePermission struct {
	Id           string `json:"id"`
	RoleId       string `json:"role_id"`
	PermissionId string `json:"permission_id"`
	Permission   []Permission
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
