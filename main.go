package main

import (
	"dh-backend-auth-sv/internal/services"
	"log"

	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//if err := godotenv.Load(); err != nil {
	//	log.Printf("Error loading app.env with godotenv: %s", err)
	//}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.AutomaticEnv()
		} else {
			log.Fatalf("cannot read config: %v", err)
		}
	}

	services.Start()
}
