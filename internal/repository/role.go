package repository

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

func (r *Repository) CreateRole(ctx context.Context, title string) error {
	role := &model.Role{
		Title: title,
	}
	return r.db.Save(role)
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

func (r *Repository) DeleteRole(ctx context.Context, id string) error {
	dest := &model.Role{}
	return r.db.Delete(dest, id)
}
