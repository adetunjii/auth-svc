package sqlstore

import (
	"context"
	"errors"

	"github.com/adetunjii/auth-svc/internal/model"
	"gorm.io/gorm"
)

var (
	ErrInvalidRoleId = errors.New("invalid role id")
	ErrRoleNotFound  = errors.New("role not found")
)

type SqlRoleStore struct {
	*SqlStore
}

func newRoleStore(sqlStore *SqlStore) *SqlRoleStore {
	return &SqlRoleStore{sqlStore}
}

func (rs SqlRoleStore) Create(ctx context.Context, arg *model.Role) error {
	return rs.db.Save(arg)
}

func (rs SqlRoleStore) List(ctx context.Context, conds map[string]interface{}, page int, size int) ([]*model.Role, error) {

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 20
	}

	limit := size
	offset := (page - 1) * size
	dest := []*model.Role{}

	err := rs.db.List(&dest, conds, limit, offset)
	return dest, err
}

func (rs SqlRoleStore) FindById(ctx context.Context, id string) (*model.Role, error) {
	dest := &model.Role{}

	err := rs.db.FindById(dest, id)
	if err != nil {
		return nil, wrappedRoleError(err)
	}

	return dest, nil
}

func (rs SqlRoleStore) FindByName(ctx context.Context, name string) (*model.Role, error) {
	dest := &model.Role{}
	condition := map[string]interface{}{
		"title": name,
	}

	err := rs.db.FindOne(dest, condition)
	if err != nil {
		return nil, wrappedRoleError(err)
	}

	return dest, nil
}

func (rs SqlRoleStore) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidRoleId
	}

	dest := &model.Role{}
	err := rs.db.Delete(dest, id)
	if err != nil {
		return wrappedRoleError(err)
	}

	return nil
}

func wrappedRoleError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRoleNotFound
	}

	return err
}
