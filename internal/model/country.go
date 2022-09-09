package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Country struct {
	Id        string    `json:"id"`
	Iso       string    `json:"iso"`
	Name      string    `json:"name"`
	PhoneCode string    `json:"phone_code"`
	NiceName  string    `json:"nice_name"`
	Currency  string    `json:"currency"`
	NumCode   string    `json:"numcode" gorm:"column:numcode"`
	ImageUrl  string    `json:"imageUrl" gorm:"column:imageUrl"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Country) BeforeCreate(tx *gorm.DB) (err error) {
	c.Id = uuid.NewString()

	c.Name = strings.ToLower(c.Name)
	c.Iso = strings.ToUpper(c.Iso)
	c.Currency = strings.ToUpper(c.Currency)

	c.CreatedAt = time.Now()

	return
}
