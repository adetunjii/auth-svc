package config

// import (
// 	"dh-backend-auth-sv/internal/db/postgres"
// 	"dh-backend-auth-sv/internal/db/rediscache"
// 	"dh-backend-auth-sv/internal/helpers"
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	DbHost        string `mapstructure:"DB_HOST"`
// 	DbPort        string `mapstructure:"DB_PORT"`
// 	DbUser        string `mapstructure:"DB_USER"`
// 	DbName        string `mapstructure:"DB_NAME"`
// 	DbPassword    string `mapstructure:"DB_PASSWORD"`
// 	DbUrl         string `mapstructure:"DB_URL"`
// 	RabbitMQHost  string `mapstructure:"RABBITMQ_HOST"`
// 	RabbitMQPort  string `mapstructure:"RABBITMQ_PORT"`
// 	RabbitMQUser  string `mapstructure:"RABBITMQ_USER"`
// 	RabbitMQPass  string `mapstructure:"RABBITMQ_PASS"`
// 	CloudAMQPUrl  string `mapstructure:"CLOUDAMQP_URL"`
// 	RedisHost     string `mapstructure:"REDIS_HOST"`
// 	RedisPort     string `mapstructure:"REDIS_PORT"`
// 	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
// 	AwsSecretID   string `mapstructure:"AWS_SECRET_ID"`
// 	AwsSecretKey  string `mapstructure:"AWS_SECRET_KEY"`
// 	AwsRegion     string `mapstructure:"AWS_REGION"`
// }

// type Service struct {
// 	DB    *postgres.PostgresDB
// 	Redis *rediscache.RedisCache
// }

// func LoadConfig() *Service {

// 	config := &Config{}

// 	services := &Service{
// 		DB:    nil,
// 		Redis: nil,
// 		// RabbitMQ: nil,
// 	}

// 	consulUrl := os.Getenv("CONSUL_URL")
// 	consulKey := os.Getenv("CONSUL_KEY")

// 	remoteViper := viper.New()
// 	remoteViper.AddRemoteProvider("consul", consulUrl, consulKey)
// 	remoteViper.SetConfigType("json")

// 	if err := remoteViper.ReadRemoteConfig(); err != nil {
// 		helpers.LogEvent("ERROR", fmt.Sprintf("%s", err))
// 		helpers.FailOnError(err, "cannot read remote config")
// 	}

// 	err := remoteViper.Unmarshal(config)
// 	if err != nil {
// 		helpers.LogEvent("ERROR", fmt.Sprintf("%s", err))
// 		helpers.FailOnError(err, "cannot load remote config")
// 	}

// 	dbConfig := &postgres.Config{
// 		Host:        config.DbHost,
// 		Port:        config.DbPort,
// 		User:        config.DbUser,
// 		Password:    config.DbPassword,
// 		Name:        config.DbName,
// 		DatabaseUrl: config.DbUrl,
// 	}

// 	db := postgres.New(dbConfig)

// 	services.DB = db

// 	// rabbitMqConfig := rabbitMQ.Config{
// 	// 	Host:     config.RabbitMQHost,
// 	// 	Port:     config.RabbitMQPort,
// 	// 	User:     config.RabbitMQUser,
// 	// 	Password: config.RabbitMQPass,
// 	// 	Url:      config.CloudAMQPUrl,
// 	// }

// 	// rabbitMq := rabbitMQ.New(rabbitMqConfig)

// 	// services.RabbitMQ = rabbitMq

// 	redisConfig := rediscache.Config{
// 		Host:     config.RedisHost,
// 		Password: config.RedisPassword,
// 		Expiry:   15 * time.Second,
// 	}

// 	redis := rediscache.New(redisConfig)

// 	services.Redis = redis
// 	return services
// }
