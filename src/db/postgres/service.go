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

func (postgresDB *PostgresDB) DeleteActivities(activities *models.Activities) error {
	return nil
}
