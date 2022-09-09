package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID, _ = uuid.NewRandom()
	fmt.Println(b.ID, b.CreatedAt)
	return
}
