package rediscache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Config struct {
	Host     string
	Db       int
	Password string
	Url      string
	Expiry   time.Duration
}

type Redis struct {
	client *redis.Client
	expiry time.Duration
}

func New(config *Config) *Redis {
	redis := &Redis{
		client: nil,
		expiry: config.Expiry,
	}

	if err := redis.GetClient(config); err != nil {
		log.Fatalf("connection to redis failed: %v", err)
	}
	return redis
}

func (r *Redis) GetClient(config *Config) error {
	var redisURL *redis.Options
	if config.Url == "" {
		redisURL = &redis.Options{
			Addr:     config.Host,
			Password: config.Password,
			DB:       config.Db,
		}
	} else {
		var err error
		redisURL, err = redis.ParseURL(config.Url)
		if err != nil {
			return err
		}
	}
	r.client = redis.NewClient(redisURL)
	//defer r.CloseConnection()

	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("error connecting to redis %v", err)
	}

	log.Println("connected to redis successfully...")
	return nil
}

func (r *Redis) CloseConnection() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *Redis) RestartConnection(config *Config) error {
	if r.client != nil {
		r.CloseConnection()
	}

	return r.GetClient(config)
}
