package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	r.Id = uuid.NewString()
	r.Title = strings.ToLower(r.Title)
	err = r.Validate()
	if err != nil {
		return
	}

	return
}

func (r *Role) Validate() error {
	if r.Title == "" {
		return errors.New("title cannot be empty")
	}

	return nil
}
