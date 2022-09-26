package sqlstore

import (
	"context"

	"github.com/adetunjii/auth-svc/internal/model"
)

type SqlPermissionStore struct {
	*SqlStore
}

func newPermisisonStore(sqlStore *SqlStore) *SqlPermissionStore {
	return &SqlPermissionStore{sqlStore}
}

func (ps SqlPermissionStore) Create(ctx context.Context, arg *model.Permission) error {
	return ps.db.Save(arg)
}

func (ps SqlPermissionStore) List(ctx context.Context, conditions map[string]interface{}, page int, size int) ([]*model.Permission, error) {

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 20
	}

	limit := size
	offset := (page - 1) * size
	dest := []*model.Permission{}

	err := ps.db.List(&dest, conditions, limit, offset)
	return dest, err
}

func (ps SqlPermissionStore) AssignToRole(ctx context.Context, permissionId string, roleId string) error {
	arg := &model.RolePermission{
		RoleId:       roleId,
		PermissionId: permissionId,
	}

	return ps.db.Save(arg)
}
func (ps SqlPermissionStore) FindById(ctx context.Context, id string) (*model.Permission, error) {
	dest := &model.Permission{}

	if err := ps.db.FindById(dest, id); err != nil {
		return nil, err
	}

	return dest, nil
}

func (ps SqlPermissionStore) RemoveFromRole(ctx context.Context, permissionId string, roleId string) error {
	model := &model.RolePermission{}
	condition := map[string]interface{}{
		"permission_id": permissionId,
		"role_id":       roleId,
	}

	return ps.db.DeleteOne(model, condition)
}

const find_by_role_id_query = `
	SELECT p.id, p.name FROM permissions as p
	JOIN role_permissions as rp ON rp.permission_id = p.id
	WHERE rp.role_id = ?
`

func (ps SqlPermissionStore) FindByRoleId(ctx context.Context, roleId string) ([]model.Permission, error) {
	dest := []model.Permission{}

	if err := ps.db.Raw(dest, find_by_role_id_query, roleId); err != nil {
		return nil, err
	}

	return dest, nil
}
