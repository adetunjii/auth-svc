package main

import (
	"fmt"
	"log"

	"github.com/adetunjii/auth-svc/config"
	grpchandler "github.com/adetunjii/auth-svc/internal/handler/grpc"
	"github.com/adetunjii/auth-svc/pkg/logging"
	"github.com/spf13/viper"
)

func main() {

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

	grpcServer := grpchandler.New(services, logger)
	grpcServer.Start(grpcPort)
}
