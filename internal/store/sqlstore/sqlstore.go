package sqlstore

import "github.com/adetunjii/auth-svc/internal/port"

type Stores struct {
	user       port.UserStore
	country    port.CountryStore
	role       port.RoleStore
	permission port.PermissionStore
}

type SqlStore struct {
	db     port.DB
	logger port.AppLogger
	stores Stores
}

var _ port.Store = (*SqlStore)(nil)

func NewSqlStore(db port.DB, logger port.AppLogger) *SqlStore {
	sqlstore := &SqlStore{
		db:     db,
		logger: logger,
	}

	sqlstore.stores.user = newUserStore(sqlstore)
	sqlstore.stores.country = newCountryStore(sqlstore)
	sqlstore.stores.role = newRoleStore(sqlstore)
	sqlstore.stores.permission = newPermisisonStore(sqlstore)
	return sqlstore
}

func (ss *SqlStore) User() port.UserStore {
	return ss.stores.user
}

func (ss *SqlStore) Role() port.RoleStore {
	return ss.stores.role
}

func (ss *SqlStore) Country() port.CountryStore {
	return ss.stores.country
}

func (ss *SqlStore) Permission() port.PermissionStore {
	return ss.stores.permission
}
