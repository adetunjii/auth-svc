package models

import "time"

type Permissions struct {
	Models
	Name string
}

//admin role, token, time and date, device(ip), userid)

type Activities struct {
	Role     string
	UserID   string
	Token    string
	Time     time.Time
	DeviceIP string
}
