package ports

import "dh-backend-auth-sv/src/models"

type DB interface {
}

type RedisCache interface {
	GetSubChannel(key string) *models.User
	SaveSubChannel(key string, channel *models.User) error
}
