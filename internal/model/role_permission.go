package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermission struct {
	Id           string       `json:"id"`
	RoleId       string       `json:"role_id"`
	PermissionId string       `json:"permission_id"`
	Permissions  []Permission `gorm:"-"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`
}

func (rp *RolePermission) BeforeCreate(tx *gorm.DB) (err error) {
	rp.Id = uuid.NewString()
	return
}
