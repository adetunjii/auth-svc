package model

import "time"

type LoginActivity struct {
	BaseEntity

	UserId string    `json:"user_id"`
	Token  string    `json:"token"`
	Device string    `json:"device"`
	Time   time.Time `json:"time"`
}
