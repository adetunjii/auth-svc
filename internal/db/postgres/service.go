package postgres

import (
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"fmt"
)

func (postgresDB *PostgresDB) SaveActivities(activities *models.Activities) error {
	helpers.LogEvent("INFO", fmt.Sprintf("saving image"))
	err := postgresDB.DB.Create(activities).Error
	return err
}

func (postgresDB *PostgresDB) DeleteActivities(userID string) error {
	helpers.LogEvent("INFO", fmt.Sprintf("deleting user activities with user_id :%s", userID))
	interest := &models.Activities{}
	err := postgresDB.DB.Where("user_id = ?", userID).Delete(interest).Error
	return err
}

func (postgresDB *PostgresDB) GetAllCountries() ([]*models.Country, error) {
	helpers.LogEvent("INFO", fmt.Sprintf("getting all countries"))

	var country []*models.Country
	err := postgresDB.DB.Find(&country).Error
	return country, err
}
