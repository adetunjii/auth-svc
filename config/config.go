package config

import (
	"fmt"
	"time"

	"gitlab.com/dh-backend/auth-service/internal/db"
	"gitlab.com/dh-backend/auth-service/internal/port"
	"gitlab.com/dh-backend/auth-service/internal/repository"
	"gitlab.com/dh-backend/auth-service/internal/services/rabbitmq"
	"gitlab.com/dh-backend/auth-service/internal/services/redis"
	"gitlab.com/dh-backend/auth-service/internal/util"
)

type Config struct {
	DbHost        string `mapstructure:"DB_HOST" json:"dbHost"`
	DbPort        string `mapstructure:"DB_PORT" json:"dbPort"`
	DbUser        string `mapstructure:"DB_USER" json:"dbUser"`
	DbName        string `mapstructure:"DB_NAME" json:"dbName"`
	DbPassword    string `mapstructure:"DB_PASSWORD" json:"dbPassword"`
	DbUrl         string `mapstructure:"DB_URL" json:"dbUrl"`
	RabbitMQHost  string `mapstructure:"RABBITMQ_HOST" json:"rabbitMQHost"`
	RabbitMQPort  string `mapstructure:"RABBITMQ_PORT" json:"rabbitMQPort"`
	RabbitMQUser  string `mapstructure:"RABBITMQ_USER" json:"rabbitMQUser"`
	RabbitMQPass  string `mapstructure:"RABBITMQ_PASS" json:"rabbitMQPass"`
	CloudAMQPUrl  string `mapstructure:"CLOUDAMQP_URL" json:"cloudAMQPUrl"`
	RedisHost     string `mapstructure:"REDIS_HOST" json:"redisHost"`
	RedisPort     string `mapstructure:"REDIS_PORT" json:"redisPort"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" json:"redisPassword"`
	AwsSecretID   string `mapstructure:"AWS_SECRET_ID" json:"awsSecretID"`
	AwsSecretKey  string `mapstructure:"AWS_SECRET_KEY" json:"awsSecretKey"`
	AwsRegion     string `mapstructure:"AWS_REGION" json:"awsRegion"`
	JwtSecretKey  string `mapstructure:"JWT_SECRETKEY" json:"jwtSecretKey"`
}

type Service struct {
	Repository *repository.Repository
	RabbitMQ   *rabbitmq.Connection
	Redis      *redis.Redis
	JwtFactory *util.JwtFactory
}

func LoadConfig(logger port.AppLogger) *Service {

	services := &Service{}

	config, err := VaultSecrets()
	if err != nil {
		logger.Fatal("failed to load vault secrets", err)
	}

	dbConfig := db.DBConfig{
		Host:        config.DbHost,
		Port:        config.DbPort,
		User:        config.DbUser,
		Name:        config.DbName,
		Password:    config.DbPassword,
		DatabaseUrl: config.DbUrl,
	}

	db := db.New(dbConfig, logger)

	repository := repository.New(db, logger)
	services.Repository = repository

	redisConfig := &redis.Config{
		Host:     config.RedisHost,
		Password: config.RedisPassword,
		Expiry:   time.Second * 15,
		// Url:      config.CloudAMQPUrl,
	}

	redis := redis.New(redisConfig, logger)

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

	rabbitMQ, err := rabbitmq.NewConnection("dh-queue", "dh-exchange", []string{"notification_queue"}, logger, rabbitMQUrl)
	if err != nil {
		logger.Fatal("failed to start rabbitmq", err)
	}

	services.RabbitMQ = rabbitMQ

	jwtFactory, err := util.NewJwtFactory(config.JwtSecretKey)
	if err != nil {
		logger.Fatal("invalid jwt secret key", err)
	}

	services.JwtFactory = jwtFactory

	return services

}
