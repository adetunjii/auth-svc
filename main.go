package main

import (
	"dh-backend-auth-sv/internal/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env with godotenv: %s", err)
	}

	// addr := os.Getenv("VAULT_ADDR")
	// secretPath := os.Getenv("VAULT_SECRET_PATH")

	// config.NewVaultClient()

	services.Start()
}
