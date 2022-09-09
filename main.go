package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gitlab.com/dh-backend/auth-service/config"
	grpcHandler "gitlab.com/dh-backend/auth-service/internal/handler/grpc"
	"gitlab.com/dh-backend/auth-service/pkg/logging"
)

func main() {

	js := "UserJwtSecretKeyHasToBe32CharactersLong"
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(js)))

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

	grpcPort := fmt.Sprintf(":%s", viper.GetString("PORT"))
	if grpcPort == ":" || grpcPort == "" {
		grpcPort = ":8080"
	}

	zapSugarLogger := logging.NewZapSugarLogger()
	logger := logging.NewLogger(zapSugarLogger)
	services := config.LoadConfig(logger)

	grpcServer := grpcHandler.New(services.Repository, services.Redis, services.RabbitMQ, services.JwtFactory, logger)
	grpcServer.Start(grpcPort)
}
