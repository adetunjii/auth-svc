package ports

import "dh-backend-auth-sv/src/models"

type DB interface {
	SaveActivities(activities *models.Activities) error
	DeleteActivities(activities *models.Activities) error
}

type RedisCache interface {
	GetSubChannel(key string) *models.User
	SaveSubChannel(key string, channel *models.User) error
	SaveRoleChannel(key string, channel []byte) error
	GetRoleChannels(key string) []models.UserRole
}
