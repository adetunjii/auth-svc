package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	p.Id = uuid.NewString()
	p.Name = strings.ToLower(p.Name)
	err = p.Validate()
	if err != nil {
		return
	}

	return
}

func (p *Permission) Validate() error {
	if p.Name == "" {
		return errors.New("permission name cannot be empty")
	}

	return nil
}
