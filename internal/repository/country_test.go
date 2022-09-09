package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/dh-backend/auth-service/internal/model"
)

func TestCreateCountry(t *testing.T) {
	country := &model.Country{
		Name:      "Nigeria",
		Iso:       "NG",
		PhoneCode: "234",
		NiceName:  "Nigeria",
		Currency:  "NGN",
		NumCode:   "566",
	}

	err := testRepo.CreateCountry(context.Background(), country)
	require.NoError(t, err)
}

func TestListCountries(t *testing.T) {
	countries, err := testRepo.ListCountries(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, countries)
}

func TestGetCountryByName(t *testing.T) {
	country, err := testRepo.GetCountryByName(context.Background(), "nigeria")
	require.NoError(t, err)
	require.NotEmpty(t, country)

	require.Equal(t, country.Name, "nigeria")
	require.Equal(t, country.Iso, "NG")
	require.Equal(t, country.NiceName, "Nigeria")
	require.Equal(t, country.Currency, "NGN")
	require.Equal(t, country.PhoneCode, "234")
}
