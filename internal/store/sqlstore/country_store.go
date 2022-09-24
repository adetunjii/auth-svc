package sqlstore

import (
	"context"
	"errors"
	"strings"

	"gitlab.com/dh-backend/auth-service/internal/model"
	"gorm.io/gorm"
)

var ErrCountryNotFound = errors.New("country not found")

type SqlCountryStore struct {
	*SqlStore
}

func newCountryStore(sqlStore *SqlStore) *SqlCountryStore {
	return &SqlCountryStore{sqlStore}
}

func (cs SqlCountryStore) Create(ctx context.Context, arg *model.Country) error {
	return cs.db.Save(arg)
}

func (cs SqlCountryStore) List(ctx context.Context) ([]*model.Country, error) {
	dest := []*model.Country{}
	err := cs.db.FindAll(&dest, nil)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (cs SqlCountryStore) GetByName(ctx context.Context, name string) (*model.Country, error) {

	country_name := strings.ToLower(name)
	country := &model.Country{}
	condition := map[string]interface{}{
		"name": country_name,
	}

	err := cs.db.FindOne(country, condition)
	if err != nil {
		return nil, wrappedCountryError(err)
	}

	return country, nil
}

func wrappedCountryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCountryNotFound
	}

	return err
}
