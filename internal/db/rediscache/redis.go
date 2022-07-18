package rediscache

import (
	"dh-backend-auth-sv/internal/ports"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

type RedisCache struct {
	host     string
	password string
	db       int
	expires  time.Duration
}

func NewRedisCache(host string, password string, db int, expires time.Duration) ports.RedisCache {
	fmt.Println(host, password, db)
	return &RedisCache{
		host:     host,
		password: password,
		db:       db,
		expires:  expires,
	}
}

func (r *RedisCache) GetClient() *redis.Client {
	var redisURL *redis.Options

	fmt.Println(r.host)

	if os.Getenv("REDIS_URL") == "" {
		redisURL = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: os.Getenv("REDIS_PASS"),
			DB:       10,
		}
	} else {
		var err error
		redisURL, err = redis.ParseURL(os.Getenv("REDIS_URL"))
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Println(redisURL.Password)
	return redis.NewClient(redisURL)
}
