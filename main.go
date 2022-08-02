package main

import (
	"dh-backend-auth-sv/internal/services"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//if err := godotenv.Load(); err != nil {
	//	log.Printf("Error loading app.env with godotenv: %s", err)
	//}

	services.Start()
}
