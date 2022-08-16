package config

import (
	"dh-backend-auth-sv/internal/db/postgres"
	"dh-backend-auth-sv/internal/db/rediscache"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/rabbitMQ"
	"fmt"
	"time"
)

type Config struct {
	DbHost         string `mapstructure:"DB_HOST" json:"dbHost"`
	DbPort         string `mapstructure:"DB_PORT" json:"dbPort"`
	DbUser         string `mapstructure:"DB_USER" json:"dbUser"`
	DbName         string `mapstructure:"DB_NAME" json:"dbName"`
	DbPassword     string `mapstructure:"DB_PASSWORD" json:"dbPassword"`
	DbUrl          string `mapstructure:"DB_URL" json:"dbUrl"`
	RabbitMQHost   string `mapstructure:"RABBITMQ_HOST" json:"rabbitMQHost"`
	RabbitMQPort   string `mapstructure:"RABBITMQ_PORT" json:"rabbitMQPort"`
	RabbitMQUser   string `mapstructure:"RABBITMQ_USER" json:"rabbitMQUser"`
	RabbitMQPass   string `mapstructure:"RABBITMQ_PASS" json:"rabbitMQPass"`
	CloudAMQPUrl   string `mapstructure:"CLOUDAMQP_URL" json:"cloudAMQPUrl"`
	RedisHost      string `mapstructure:"REDIS_HOST" json:"redisHost"`
	RedisPort      string `mapstructure:"REDIS_PORT" json:"redisPort"`
	RedisPassword  string `mapstructure:"REDIS_PASSWORD" json:"redisPassword"`
	AwsSecretID    string `mapstructure:"AWS_SECRET_ID" json:"awsSecretID"`
	AwsSecretKey   string `mapstructure:"AWS_SECRET_KEY" json:"awsSecretKey"`
	AwsRegion      string `mapstructure:"AWS_REGION" json:"awsRegion"`
	UserServiceUrl string `mapstructure:"USER_SERVICE_URL" json:"userServiceUrl"`
}

type Service struct {
	DB       *postgres.PostgresDB
	Redis    *rediscache.Redis
	RabbitMQ *rabbitMQ.RabbitMQ
}

func LoadConfig() *Service {

	services := &Service{}

	config, err := VaultSecrets()
	if err != nil {
		helpers.LogEvent("ERROR", "couldn't load secrets")
		helpers.FailOnError(err, "couldn't load secrets")
	}

	dbConfig := &postgres.Config{
		Host:        config.DbHost,
		Port:        config.DbPort,
		User:        config.DbUser,
		Name:        config.DbName,
		Password:    config.DbPassword,
		DatabaseUrl: config.DbUrl,
	}

	db := postgres.New(dbConfig)

	services.DB = db

	redisConfig := &rediscache.Config{
		Host:     config.RedisHost,
		Password: config.RedisPassword,
		Expiry:   time.Second * 15,
	}

	redis := rediscache.New(redisConfig)

	services.Redis = redis

	var rabbitMQUrl string
	if config.CloudAMQPUrl == "" {
		rabbitMQUrl = fmt.Sprintf(
			"amqp://%s:%s@%s:%s",
			config.RabbitMQUser,
			config.RabbitMQPass,
			config.RabbitMQHost,
			config.RabbitMQPort,
		)

	} else {
		rabbitMQUrl = config.CloudAMQPUrl
	}

	rabbitMQ := rabbitMQ.New(rabbitMQUrl)

	services.RabbitMQ = rabbitMQ

	return services
}
