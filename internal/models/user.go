package models

import "time"

type User struct {
	Models
	RoleID          string `json:"roleID"`
	Role            Role   `gorm:"foreignkey:RoleID" protobuf:"bytes,1,opt,name=roles" json:"roles"`
	Email           string `protobuf:"bytes,2,opt,name=email" gorm:"unique"`
	PhoneNumber     string `protobuf:"bytes,3,opt,name=phoneNumber" gorm:"unique"`
	PhoneCode       string `json:"phoneCode"`
	HashedPassword  string `protobuf:"bytes,4,opt,name=hashedPassword"`
	IsActive        bool   `json:"is_active" gorm:"default:false" protobuf:"bytes,5,opt,name=isActive"`
	IsEmailVerified bool   `json:"is_email_verified,omitempty" gorm:"default:false"`
	IsPhoneVerified bool   `json:"is_phone_verified,omitempty" gorm:"default:false"`
	UserInformation `protobuf:"bytes,6,opt,name=userInformation"`
	UserInterests   []*Interest     `gorm:"many2many:user_interests"protobuf:"bytes,7,opt,name=userInterests"`
	Picture         []*ProfileImage `gorm:"many2many:user_pictures"protobuf:"bytes,8,opt,name=picture"`
}

type Interest struct {
	Models
	Item string `json:"item"`
}

type UserInterests struct {
	UserID     string
	InterestID string
}

type Permission struct {
	Models
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePermissions struct {
	RoleID       string
	PermissionID string
}

type ProfileImage struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImageType string    `json:"image_type"`
	Path      string    `json:"path"`
}
