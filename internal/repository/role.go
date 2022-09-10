package repository

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

func (r *Repository) CreateRole(ctx context.Context, arg *model.Role) error {
	return r.db.Save(arg)
}

func (r *Repository) FindRoleById(ctx context.Context, id string) (*model.Role, error) {
	dest := &model.Role{}

	err := r.db.FindById(dest, id)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *Repository) ListRoles(context context.Context) ([]*model.Role, error) {
	dest := []*model.Role{}

	err := r.db.FindAll(&dest, nil)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
