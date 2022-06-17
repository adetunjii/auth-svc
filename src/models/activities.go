package models

import "time"

type Activities struct {
	ID       string
	UserID   string
	Token    string
	Time     time.Time
	DeviceIP int
}

type ActivityRoles struct {
	Models
	RoleName string `json:"role_name"bson:"role_name"`
}
