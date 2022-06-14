package models

import "time"

type ProfileImage struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImageType string    `json:"image_type"`
	Path      string    `json:"path"`
}
