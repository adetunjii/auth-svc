package port

import (
	"context"

	"github.com/adetunjii/auth-svc/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, arg *model.User) error
	ListUsers(ctx context.Context, conds map[string]interface{}, page int, size int) ([]*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserById(ctx context.Context, id string) (*model.User, error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string, phoneCode string) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, arg *model.User) error
	CreateRole(ctx context.Context, title string) error
	ListRoles(ctx context.Context) ([]*model.Role, error)
	DeleteRole(ctx context.Context, id string) error
	FindRoleById(ctx context.Context, id string) (*model.Role, error)
	CreatePermission(ctx context.Context, arg *model.Permission) error
	AssignPermissionToRole(ctx context.Context, roleId string, permissionId string) error
	FindPermissionsByRoleId(ctx context.Context, roleId string) ([]*model.RolePermission, error)
	CreateInterest(ctx context.Context, arg *model.Interest) error
	CreateUserInterest(ctx context.Context, userId string, interestId string) error
	CreateCountry(ctx context.Context, arg *model.Country) error
	ListCountries(ctx context.Context) ([]*model.Country, error)
	GetCountryByName(ctx context.Context, name string) (*model.Country, error)
}
