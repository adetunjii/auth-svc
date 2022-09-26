package model

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/adetunjii/auth-svc/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	UserPasswordMinLength = 8
	UserPasswordMaxLength = 72
	UserEmailMaxLength    = 128
)

type User struct {
	Id                 string    `json:"id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Username           string    `json:"username,omitempty"`
	Email              string    `json:"email"`
	PhoneNumber        string    `json:"phone_number"`
	PhoneCode          string    `json:"phone_code"`
	RoleId             string    `json:"role_id,omitempty"`
	Password           string    `json:"password"`
	LastPasswordUpdate time.Time `json:"last_password_update,omitempty"`
	IsEmailVerified    bool      `json:"is_email_verified,omitempty"`
	IsPhoneVerified    bool      `json:"is_phone_verified,omitempty"`
	IsActive           bool      `json:"is_active,omitempty"`
	Address            string    `json:"address"`
	State              string    `json:"state,omitempty"`
	Country            string    `json:"country"`
	Timezone           string    `json:"timezone,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// user update struct
type UserPatch struct {
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	Username        *string `json:"username"`
	Password        *string `json:"password,omitempty"`
	IsEmailVerified *bool   `json:"is_email_verified"`
	IsPhoneVerified *bool   `json:"is_phone_verified"`
	IsActive        *bool   `json:"is_active"`
	Address         *string `json:"address"`
	State           *string `json:"state"`
	Country         *string `json:"country"`
	Timezone        *string `json:"timezone"`
}

// TODO: validate country: check if user is from a supported country
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.NewString()

	if u.Username == "" {
		u.Username = u.FirstName
	}

	// set default role to user
	if u.RoleId == "" {
		role := &Role{}
		err = tx.Where("title = ?", "user").Find(role).Error
		if err != nil {
			return
		}

		u.RoleId = role.Id
	}

	u.Username = SanitizeUnicode(u.Username)
	u.FirstName = SanitizeUnicode(u.FirstName)
	u.LastName = SanitizeUnicode(u.LastName)
	u.Address = SanitizeUnicode(u.Address)
	u.State = SanitizeUnicode(u.State)

	u.Username = strings.ToLower(u.Username)
	u.FirstName = strings.ToLower(u.FirstName)
	u.LastName = strings.ToLower(u.LastName)
	u.Address = strings.ToLower(u.Address)
	u.State = strings.ToLower(u.State)
	u.Address = strings.ToLower(u.Address)
	u.Email = strings.ToLower(u.Email)

	err = u.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	u.PhoneNumber = TrimPhoneNumber(u.PhoneNumber, u.PhoneCode)

	// hash password
	if u.Password != "" {
		u.Password, err = HashPassword(u.Password)
		if err != nil {
			return
		}
	}

	u.CreatedAt = time.Now()

	return nil
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {

	u.Username = SanitizeUnicode(u.Username)
	u.FirstName = SanitizeUnicode(u.FirstName)
	u.LastName = SanitizeUnicode(u.LastName)
	u.Address = SanitizeUnicode(u.Address)
	u.State = SanitizeUnicode(u.State)

	u.Username = strings.ToLower(u.Username)
	u.FirstName = strings.ToLower(u.FirstName)
	u.LastName = strings.ToLower(u.LastName)
	u.Address = strings.ToLower(u.Address)
	u.State = strings.ToLower(u.State)
	u.Address = strings.ToLower(u.Address)
	u.Email = strings.ToLower(u.Email)

	if u.Password != "" {
		u.Password, err = HashPassword(u.Password)
		if err != nil {
			return
		}
	}

	u.UpdatedAt = time.Now()

	return
}

// func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
// 	fmt.Println(u.Id)
// 	hashPassword, err := HashPassword(u.Password)
// 	if err != nil {
// 		return
// 	}

// 	tx.Model(&User{}).Where("id = ?", u.Id).Update("password", hashPassword)

// 	return
// }

// remove the first zero for Nigerian numbers
func TrimPhoneNumber(phoneNumber string, phoneCode string) string {
	if phoneCode == "234" && phoneNumber[0] == '0' {
		return phoneNumber[1:]
	}
	return phoneNumber
}

func (u *User) Validate() error {
	if u.Id == "" {
		return ErrInvalidUserId
	}

	if err := IsValidId(u.Id); err != nil {
		return err
	}

	if err := IsValidEmail(u.Email); err != nil {
		return err
	}

	if u.Password != "" {
		if err := IsPasswordValid(u.Password); err != nil {
			return err
		}
	}

	if u.FirstName == "" {
		return errors.New("first name cannot be empty")
	}

	if u.LastName == "" {
		return errors.New("last name cannot be empty")
	}

	if u.PhoneCode == "" || len(u.PhoneCode) == 0 {
		return errors.New("phone code cannot be empty")
	}

	if u.PhoneNumber == "" || len(u.PhoneNumber) == 0 {
		return errors.New("phone number cannot be empty")
	}

	if u.Country == "" {
		return errors.New("country cannot be empty")
	}

	return nil
}

func (u *User) Patch(patch *UserPatch) (err error) {
	if patch.FirstName != nil {
		u.FirstName = *patch.FirstName
	}

	if patch.LastName != nil {
		u.LastName = *patch.LastName
	}

	if patch.Username != nil {
		u.Username = *patch.Username
	}

	if patch.Password != nil {
		u.Password, err = HashPassword(*patch.Password)
		if err != nil {
			return
		}
	}

	if patch.IsEmailVerified != nil {
		u.IsEmailVerified = *patch.IsEmailVerified
	}

	if patch.IsPhoneVerified != nil {
		u.IsPhoneVerified = *patch.IsPhoneVerified
	}

	if patch.IsActive != nil {
		u.IsActive = *patch.IsActive
	}

	if patch.Address != nil {
		u.Address = *patch.Address
	}

	if patch.State != nil {
		u.State = *patch.State
	}

	if patch.Country != nil {
		u.Country = *patch.Country
	}

	if patch.Timezone != nil {
		u.Timezone = *patch.Timezone
	}

	return
}

func IsPasswordValid(password string) error {

	if len(password) < UserPasswordMinLength {
		return ErrInvalidPasswordTooLong
	}

	if len(password) > UserPasswordMaxLength {
		return ErrInvalidPasswordTooLong
	}

	if !strings.ContainsAny(password, util.LowercaseAlphabets) ||
		!strings.ContainsAny(password, util.UppercaseAlphabets) ||
		!strings.ContainsAny(password, util.Symbols) ||
		!strings.ContainsAny(password, util.Digits) {
		return ErrInvalidPasswordPattern
	}

	return nil
}

// check if id is a valid uuid [ex: 12b33818-626e4-8808-388c-598856b672c (8-4-4-4-12)]
func IsValidId(value string) error {
	if len(value) != 36 {
		return ErrInvalidUserId
	}

	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && string(r) != "-" {
			return ErrInvalidUserId
		}
	}
	return nil
}

func IsValidEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmailAddress
	}

	if addr.Name != "" {
		return ErrInvalidEmailAddress
	}

	return nil
}
