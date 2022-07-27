package postgres

import (
	"dh-backend-auth-sv/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Config struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	DatabaseUrl string `json:"url"`
}

type PostgresDB struct {
	DB *gorm.DB
}

func New(config *Config) *PostgresDB {
	db := &PostgresDB{
		DB: nil,
	}

	if err := db.Connect(config); err != nil {
		log.Fatalf("connection to db failed: %v", err)
	}
	return db
}

func (postgresDB *PostgresDB) Connect(config *Config) error {
	fmt.Println(config)

	var dns string
	databaseUrl := config.DatabaseUrl
	if databaseUrl == "" {
		dns = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.Name)
	} else {
		dns = databaseUrl
	}
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}

	postgresDB.DB = db

	//Migrate the schema
	//err = postgresDB.DB.AutoMigrate(&models.Role{}, &models.Interest{}, &models.User{}, &models.UserRole{}, &models.Country{})
	err = postgresDB.DB.AutoMigrate(&models.Country{})
	if err != nil {
		return err
	}

	log.Println("Database Connected Successfully...")
	return nil
}

func (postgresDB *PostgresDB) CloseConnection() error {
	if postgresDB.DB != nil {
		connection, err := postgresDB.DB.DB()
		if err != nil {
			return err
		}

		connection.Close()
	}
	return nil
}

func (postgresDB *PostgresDB) RestartConnection(config *Config) error {
	if postgresDB.DB != nil {
		postgresDB.CloseConnection()
	}

	return postgresDB.Connect(config)
}
