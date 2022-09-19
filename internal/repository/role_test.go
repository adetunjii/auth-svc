package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
)

func TestCreateRole(t *testing.T) {

	err := testRepo.CreateRole(context.Background(), "user")
	require.NoError(t, err)
}

func TestDuplicateRoleError(t *testing.T) {
	err := testRepo.CreateRole(context.Background(), "user")
	require.Error(t, err)

	pgErr, ok := err.(*pgconn.PgError)
	require.True(t, ok)
	require.NotEmpty(t, pgErr)
	require.Equal(t, pgErr.Code, "23505")

}

func TestGetAllRoles(t *testing.T) {
	roles, err := testRepo.ListRoles(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, roles)
}
