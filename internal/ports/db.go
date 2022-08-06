package ports

import (
	"dh-backend-auth-sv/internal/models"
)

type DB interface {
	SaveActivities(activities *models.Activities) error
	DeleteActivities(userID string) error
	GetAllCountries() ([]*models.Country, error)
}

type RedisCache interface {
	GetSubChannel(key string) *models.User
	SaveSubChannel(key string, channel *models.User) error
	SaveRoleChannel(key string, channel []byte) error
	GetRoleChannels(key string) []models.UserRole
	SaveOTP(key string, otpType string, value any) error
	GetOTP(key string, otpType string) (*models.OtpVerification, error)
}
