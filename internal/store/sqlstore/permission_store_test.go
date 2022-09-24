package sqlstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/dh-backend/auth-service/internal/model"
)

var testPermission = &model.Permission{
	Name:        "delete user",
	Description: "Test permission",
}

func TestCreatePermission(t *testing.T) {
	err := sqlStore.Permission().Create(context.Background(), testPermission)
	require.NoError(t, err)
}

func TestAssignToRole(t *testing.T) {
	roleId := "9fbb4d24-26f1-4e1c-a699-2d9a401371e9"
	permissionId := "2fbaa395-c854-40f2-92a1-77a3620ecf99"
	permission, err := sqlStore.Permission().FindById(context.Background(), permissionId)
	require.NoError(t, err)
	require.NotEmpty(t, permission)

	err = sqlStore.Permission().AssignToRole(context.Background(), permissionId, roleId)
	require.NoError(t, err)
}

func TestRemovePermissionFromRole(t *testing.T) {

	roleId := "9fbb4d24-26f1-4e1c-a699-2d9a401371e9"
	permissionId := "2fbaa395-c854-40f2-92a1-77a3620ecf99"

	err := sqlStore.Permission().RemoveFromRole(context.Background(), permissionId, roleId)
	require.NoError(t, err)
}
