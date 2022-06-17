package models

import "time"

type Activities struct {
	ID     string
	UserID string
	Token  string
	Time   time.Time
	Device string
}

//type Device struct {
//	IP string
//}

type ActivityRoles struct {
	Models
	RoleName string `json:"role_name"bson:"role_name"`
}
