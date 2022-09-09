package repository

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

func (r *Repository) CreateUser(ctx context.Context, arg *model.User) error {
	return r.db.Save(arg)
}

func (r *Repository) ListUsers(ctx context.Context, conds map[string]interface{}, page int, size int) ([]*model.User, error) {

	limit := size
	offset := (page - 1) * size
	dest := []*model.User{}

	err := r.db.List(&dest, conds, limit, offset)
	return dest, err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	dest := &model.User{}
	conds := map[string]interface{}{
		"email": email,
	}

	err := r.db.FindOne(dest, conds)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *Repository) FindUserByPhoneNumber(ctx context.Context, phoneNumber string, phoneCode string) (*model.User, error) {
	dest := &model.User{}
	conds := map[string]interface{}{
		"phone_number": phoneNumber,
		"phone_code":   phoneCode,
	}

	err := r.db.FindOne(dest, conds)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *Repository) FindUserById(ctx context.Context, id string) (*model.User, error) {
	dest := &model.User{}

	err := r.db.FindById(dest, id)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id string) error {
	dest := &model.User{}
	return r.db.Delete(dest, id)
}

func (r *Repository) UpdateUser(ctx context.Context, id string, arg *model.User) error {
	dest := &model.User{}
	conds := map[string]interface{}{
		"id": id,
	}

	return r.db.Update(dest, conds, arg)
}
