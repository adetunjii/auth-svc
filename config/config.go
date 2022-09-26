package config

import (
	"fmt"
	"time"

	"github.com/adetunjii/auth-svc/internal/db"
	"github.com/adetunjii/auth-svc/internal/port"
	"github.com/adetunjii/auth-svc/internal/services/oauth"
	"github.com/adetunjii/auth-svc/internal/services/rabbitmq"
	"github.com/adetunjii/auth-svc/internal/services/redis"
	"github.com/adetunjii/auth-svc/internal/store/sqlstore"
	"github.com/adetunjii/auth-svc/internal/util"
)

type Config struct {
	DbHost             string   `mapstructure:"DB_HOST" json:"dbHost"`
	DbPort             string   `mapstructure:"DB_PORT" json:"dbPort"`
	DbUser             string   `mapstructure:"DB_USER" json:"dbUser"`
	DbName             string   `mapstructure:"DB_NAME" json:"dbName"`
	DbPassword         string   `mapstructure:"DB_PASSWORD" json:"dbPassword"`
	DbUrl              string   `mapstructure:"DB_URL" json:"dbUrl"`
	RabbitMQHost       string   `mapstructure:"RABBITMQ_HOST" json:"rabbitMQHost"`
	RabbitMQPort       string   `mapstructure:"RABBITMQ_PORT" json:"rabbitMQPort"`
	RabbitMQUser       string   `mapstructure:"RABBITMQ_USER" json:"rabbitMQUser"`
	RabbitMQPass       string   `mapstructure:"RABBITMQ_PASS" json:"rabbitMQPass"`
	CloudAMQPUrl       string   `mapstructure:"CLOUDAMQP_URL" json:"cloudAMQPUrl"`
	RedisHost          string   `mapstructure:"REDIS_HOST" json:"redisHost"`
	RedisPort          string   `mapstructure:"REDIS_PORT" json:"redisPort"`
	RedisPassword      string   `mapstructure:"REDIS_PASSWORD" json:"redisPassword"`
	JwtSecretKey       string   `mapstructure:"JWT_SECRETKEY" json:"jwtSecretKey"`
	GoogleClientId     string   `mapstructure:"GOOGLE_CLIENT_ID" json:"googleClientId"`
	GoogleClientSecret string   `mapstructure:"GOOGLE_CLIENT_SECRET" json:"googleClientSecret"`
	GoogleUserScopes   []string `mapstructure:"GOOGLE_USER_SCOPES" json:"googleUserScopes"`
	GoogleRedirectURL  string   `mapstructure:"GOOGLE_REDIRECT_URL" json:"googleRedirectURL"`
}

type Service struct {
	RabbitMQ     *rabbitmq.Connection
	Redis        *redis.Redis
	JwtFactory   *util.JwtFactory
	GoogleClient *oauth.GoogleClient
	Store        port.Store
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

	googleClient := oauth.NewGoogleClient(config.GoogleClientId, config.GoogleClientSecret, config.GoogleUserScopes, config.GoogleRedirectURL, logger)
	services.GoogleClient = googleClient

	sqlStore := sqlstore.NewSqlStore(db, logger)
	services.Store = sqlStore

	return services

}
