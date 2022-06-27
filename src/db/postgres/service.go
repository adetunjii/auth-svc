package postgres

import (
	"dh-backend-auth-sv/src/helpers"
	"dh-backend-auth-sv/src/models"
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
