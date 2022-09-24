package sqlstore

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/dh-backend/auth-service/internal/model"
	"gitlab.com/dh-backend/auth-service/internal/util"
)

func createTestAccount(t *testing.T) *model.User {
	user := &model.User{
		FirstName:   util.RandomName(),
		LastName:    util.RandomName(),
		Email:       fmt.Sprintf("%s@example.com", util.RandomString(6)),
		PhoneNumber: util.RandomPhoneNumber(),
		PhoneCode:   util.RandomCountryCode(),
		Password:    util.RandomPassword(),
		Address:     "Lekki phase 1",
		State:       "Lagos",
		Country:     "Nigeria",
	}

	err := sqlStore.User().Save(context.Background(), user)
	require.NoError(t, err)

	return user
}

func TestSave(t *testing.T) {
	createTestAccount(t)
}

func TestFindUserById(t *testing.T) {
	user := createTestAccount(t)

	dbUser, err := sqlStore.User().FindById(context.Background(), user.Id)
	require.NoError(t, err)

	require.Equal(t, user.Id, dbUser.Id)
	require.Equal(t, user.FirstName, dbUser.FirstName)
	require.Equal(t, user.LastName, dbUser.LastName)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.PhoneNumber, dbUser.PhoneNumber)
	require.Equal(t, user.PhoneCode, dbUser.PhoneCode)
	require.Equal(t, user.RoleId, dbUser.RoleId)
	require.Equal(t, user.Address, dbUser.Address)
	require.Equal(t, user.Country, dbUser.Country)

	noUser, err := sqlStore.User().FindById(context.Background(), util.RandomUUID())
	require.Error(t, err)
	require.Empty(t, noUser)
}

func TestUserByEmail(t *testing.T) {
	user := createTestAccount(t)

	dbUser, err := sqlStore.User().FindByEmail(context.Background(), user.Email)
	require.NoError(t, err)

	require.Equal(t, user.Id, dbUser.Id)
	require.Equal(t, user.FirstName, dbUser.FirstName)
	require.Equal(t, user.LastName, dbUser.LastName)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.PhoneNumber, dbUser.PhoneNumber)
	require.Equal(t, user.PhoneCode, dbUser.PhoneCode)
	require.Equal(t, user.RoleId, dbUser.RoleId)
	require.Equal(t, user.Address, dbUser.Address)
	require.Equal(t, user.Country, dbUser.Country)

	randomEmail := fmt.Sprintf("%s@example.com", util.RandomString(6))
	noUser, err := sqlStore.User().FindByEmail(context.Background(), randomEmail)
	require.Error(t, err)
	require.Empty(t, noUser)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 5; i++ {
		createTestAccount(t)
	}

	page := 1
	size := 5

	users, err := sqlStore.User().List(context.Background(), nil, page, size)
	require.NoError(t, err)

	require.Len(t, users, 5)
}

func TestDeleteUser(t *testing.T) {
	user := createTestAccount(t)

	err := sqlStore.User().Delete(context.Background(), user.Id)
	require.NoError(t, err)

	dbUser, err := sqlStore.User().FindById(context.Background(), user.Id)
	require.Error(t, err)
	require.Empty(t, dbUser)
}

func TestUpdateUser(t *testing.T) {
	user := createTestAccount(t)

	updates := &model.User{
		FirstName: "hello",
		LastName:  "yayy",
		Username:  "teej4y",
	}

	err := sqlStore.User().Update(context.Background(), user.Id, updates)
	require.NoError(t, err)

	updatedUser, err := sqlStore.User().FindById(context.Background(), user.Id)
	require.NoError(t, err)

	require.Equal(t, updates.FirstName, updatedUser.FirstName)
	require.Equal(t, updates.LastName, updatedUser.LastName)
	require.Equal(t, updates.Username, updatedUser.Username)
}

func TestFindUserByPhoneNumber(t *testing.T) {
	user := createTestAccount(t)

	u, err := sqlStore.User().FindByPhoneNumber(context.Background(), user.PhoneNumber, user.PhoneCode)
	require.NoError(t, err)
	require.NotEmpty(t, u)

	require.Equal(t, user.Id, u.Id)
	require.Equal(t, user.FirstName, u.FirstName)
	require.Equal(t, user.LastName, u.LastName)
	require.Equal(t, user.Email, u.Email)
	require.Equal(t, user.PhoneNumber, u.PhoneNumber)
	require.Equal(t, user.PhoneCode, u.PhoneCode)
	require.Equal(t, user.RoleId, u.RoleId)
	require.Equal(t, user.Address, u.Address)
	require.Equal(t, user.Country, u.Country)
}
