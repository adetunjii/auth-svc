package repository

import (
	"context"
	"strings"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

func (r *Repository) CreateCountry(ctx context.Context, arg *model.Country) error {
	return r.db.Save(arg)
}

func (r *Repository) ListCountries(ctx context.Context) ([]*model.Country, error) {
	dest := []*model.Country{}
	err := r.db.FindAll(&dest, nil)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *Repository) GetCountryByName(ctx context.Context, name string) (*model.Country, error) {
	n := strings.ToLower(name)

	country := &model.Country{}

	condition := map[string]interface{}{
		"name": n,
	}

	err := r.db.FindOne(country, condition)
	if err != nil {
		return nil, err
	}

	return country, nil
}
