package models

type Role struct {
	Models
	Title       string
	Permissions []*Permission `gorm:"many2many:role_permissions"`
}

type UserRole struct {
	UserID string
	RoleID string
}
