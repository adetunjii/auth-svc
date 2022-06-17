package postgres

import (
	"dh-backend-auth-sv/src/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type PostgresDB struct {
	DB *gorm.DB
}

func (postgresDB *PostgresDB) Init() {
	// using gorm to connect to db driver
	var dns string
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		dns = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		fmt.Println(dns)
	} else {
		dns = databaseUrl
	}
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	postgresDB.DB = db

	err = db.AutoMigrate(models.ActivityRoles{}, &models.Activities{})
	if err != nil {
		log.Printf("Error %s", err)
	}

	if err != nil {
		log.Println("unable to create role.", err.Error())
	}
	log.Println("Database Connected Successfully...")
}
