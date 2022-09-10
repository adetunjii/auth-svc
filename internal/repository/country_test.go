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

	c := []*model.Country{
		{
			Name:      "Nigeria",
			Iso:       "NG",
			PhoneCode: "234",
			NiceName:  "Nigeria",
			Currency:  "NGN",
			NumCode:   "566",
		},
		{
			Name:      "Ghana",
			Iso:       "GH",
			PhoneCode: "233",
			NiceName:  "Ghana",
			Currency:  "GHS",
			NumCode:   "288",
		},
		{
			Name:      "Cameroon",
			Iso:       "CM",
			PhoneCode: "237",
			NiceName:  "Cameroon",
			Currency:  "CFA",
			NumCode:   "120",
		},
		{
			Name:      "Kenya",
			Iso:       "KE",
			PhoneCode: "254",
			NiceName:  "Kenya",
			Currency:  "KES",
			NumCode:   "404",
		},
		// {
		// 	Name:      "Canada",
		// 	Iso:       "CA",
		// 	PhoneCode: "1",
		// 	NiceName:  "Canada",
		// 	Currency:  "CAD",
		// 	NumCode:   "124",
		// },
		{
			Name:      "United States",
			Iso:       "US",
			PhoneCode: "1",
			NiceName:  "United States",
			Currency:  "USD",
			NumCode:   "840",
		},
		{
			Name:      "United Kingdom",
			Iso:       "UK",
			PhoneCode: "44",
			NiceName:  "United Kingdom",
			Currency:  "GBP",
			NumCode:   "826",
		},
	}

	for i := 0; i < len(c); i++ {
		err := testRepo.CreateCountry(context.Background(), c[i])
		require.NoError(t, err)
	}

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
