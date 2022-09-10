package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/dh-backend/auth-service/internal/model"
)

func TestCreateRole(t *testing.T) {
	arg := &model.Role{
		Title: "user",
	}

	err := testRepo.db.Save(arg)
	require.NoError(t, err)
}

func TestGetAllRoles(t *testing.T) {
	roles, err := testRepo.ListRoles(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, roles)
}
