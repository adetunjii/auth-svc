package repository

import (
	"gitlab.com/dh-backend/auth-service/internal/port"
)

type Repository struct {
	db     port.DB
	logger port.AppLogger
}

var _ port.Repository = (*Repository)(nil)

func New(db port.DB, logger port.AppLogger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}
