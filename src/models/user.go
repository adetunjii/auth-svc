package models

type User struct {
	Models
	Roles           []*Role `gorm:"many2many:user_roles;" protobuf:"bytes,1,opt,name=roles"`
	Email           string  `protobuf:"bytes,2,opt,name=email"`
	PhoneNumber     string  `protobuf:"bytes,3,opt,name=phoneNumber"`
	HashedPassword  string  `protobuf:"bytes,4,opt,name=hashedPassword"`
	IsActive        bool    `json:"is_active" gorm:"default:false" protobuf:"bytes,5,opt,name=isActive"`
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
