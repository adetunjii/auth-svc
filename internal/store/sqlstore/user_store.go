package sqlstore

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

type SqlUserStore struct {
	*SqlStore
}

func newUserStore(sqlstore *SqlStore) *SqlUserStore {
	return &SqlUserStore{sqlstore}
}

func (us SqlUserStore) Save(ctx context.Context, arg *model.User) error {
	return us.db.Save(arg)
}

func (us SqlUserStore) List(ctx context.Context, conds map[string]interface{}, page int, size int) ([]*model.User, error) {

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 20
	}

	limit := size
	offset := (page - 1) * size
	dest := []*model.User{}

	err := us.db.List(&dest, conds, limit, offset)
	return dest, err
}

func (us SqlUserStore) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	dest := &model.User{}
	conds := map[string]interface{}{
		"email": email,
	}

	err := us.db.FindOne(dest, conds)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (us SqlUserStore) FindById(ctx context.Context, id string) (*model.User, error) {
	dest := &model.User{}

	err := us.db.FindById(dest, id)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (us SqlUserStore) FindByPhoneNumber(ctx context.Context, phoneNumber string, phoneCode string) (*model.User, error) {
	dest := &model.User{}
	conds := map[string]interface{}{
		"phone_number": phoneNumber,
		"phone_code":   phoneCode,
	}

	err := us.db.FindOne(dest, conds)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (us SqlUserStore) Delete(ctx context.Context, id string) error {
	dest := &model.User{}
	return us.db.Delete(dest, id)
}

func (us SqlUserStore) Update(ctx context.Context, id string, arg *model.User) error {
	dest := &model.User{}
	conds := map[string]interface{}{
		"id": id,
	}

	return us.db.Update(dest, conds, arg)
}

func (us SqlUserStore) UpdatePassword(ctx context.Context, id string, password string) error {
	conds := map[string]interface{}{
		"id": id,
	}

	updates := map[string]interface{}{
		"password": password,
	}

	return us.db.Update(&model.User{}, conds, updates)
}
