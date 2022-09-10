package repository

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

// TODO: assign permission to role
func (r *Repository) AssignPermissionToRole(ctx context.Context, roleId string, permissionId string) error {
	return nil
}

func (r *Repository) FindPermissionsByRoleId(ctx context.Context, roleId string) ([]*model.RolePermission, error) {
	rp := []*model.RolePermission{}
	conditions := map[string]interface{}{
		"role_id": roleId,
	}
	err := r.db.FindWithPreload(rp, conditions, "Permission")
	if err != nil {
		return nil, err
	}

	return rp, nil
}
