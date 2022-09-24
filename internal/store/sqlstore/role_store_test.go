package sqlstore

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
	"gitlab.com/dh-backend/auth-service/internal/model"
)

var testRole = &model.Role{
	Title: "admin",
}

func TestCreateRole(t *testing.T) {

	err := sqlStore.Role().Create(context.Background(), testRole)
	require.NoError(t, err)
}

func TestDuplicateRoleError(t *testing.T) {
	err := sqlStore.Role().Create(context.Background(), testRole)
	require.Error(t, err)

	pgErr, ok := err.(*pgconn.PgError)
	require.True(t, ok)
	require.NotEmpty(t, pgErr)
	require.Equal(t, pgErr.Code, "23505")

}

func TestGetAllRoles(t *testing.T) {
	roles, err := sqlStore.Role().List(context.Background(), nil, 1, 10)
	require.NoError(t, err)
	require.NotEmpty(t, roles)
}
