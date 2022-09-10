package repository

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

func (r *Repository) CreatePermission(ctx context.Context, arg *model.Permission) error {
	return r.db.Save(arg)
}

func (r *Repository) GetAllPermissions(ctx context.Context) ([]model.Permission, error) {
	dest := []model.Permission{}

	err := r.db.FindAll(dest, nil)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
