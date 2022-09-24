package port

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

type Store interface {
	User() UserStore
	Role() RoleStore
	Permission() PermissionStore
	Country() CountryStore
	// Interest() InterestStore
}

type UserStore interface {
	Save(ctx context.Context, arg *model.User) error
	List(ctx context.Context, conds map[string]interface{}, page int, size int) ([]*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindById(ctx context.Context, id string) (*model.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string, phoneCode string) (*model.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, arg *model.User) error
	UpdatePassword(ctx context.Context, id string, password string) error
}

type RoleStore interface {
	Create(ctx context.Context, arg *model.Role) error
	List(ctx context.Context, conds map[string]interface{}, page int, size int) ([]*model.Role, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (*model.Role, error)
	FindByName(ctx context.Context, name string) (*model.Role, error)
}

type PermissionStore interface {
	Create(ctx context.Context, arg *model.Permission) error
	AssignToRole(ctx context.Context, permissionId string, roleId string) error
	FindById(ctx context.Context, id string) (*model.Permission, error)
	FindByRoleId(ctx context.Context, roleId string) ([]model.Permission, error)
	RemoveFromRole(ctx context.Context, permissionId string, roleId string) error
}

type CountryStore interface {
	Create(ctx context.Context, arg *model.Country) error
	List(ctx context.Context) ([]*model.Country, error)
	GetByName(ctx context.Context, name string) (*model.Country, error)
}

type InterestStore interface {
	Create(ctx context.Context, arg *model.Interest) error
	MatchUser(ctx context.Context, interestId string, userId string) error
}
