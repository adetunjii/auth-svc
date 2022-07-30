package main

import (
	"dh-backend-auth-sv/internal/services"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//if err := godotenv.Load(); err != nil {
	//	log.Printf("Error loading .env with godotenv: %s", err)
	//}

	viper.ReadInConfig()
	viper.AutomaticEnv()

	port := viper.Get("PORT")
	fmt.Println(port)

	services.Start()
}
