package rediscache

import (
	"dh-backend-auth-sv/src/ports"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, expires time.Duration) ports.RedisCache {
	return &RedisCache{
		host:    host,
		db:      db,
		expires: expires,
	}
}

func (r *RedisCache) getClient() *redis.Client {
	var redisURL *redis.Options
	if os.Getenv("REDIS_URL") == "" {
		redisURL = &redis.Options{
			Addr:     r.host,
			Password: "",
			DB:       r.db,
		}
	} else {
		var err error
		redisURL, err = redis.ParseURL(os.Getenv("REDIS_URL"))
		if err != nil {
			log.Println(err)
		}
	}
	return redis.NewClient(redisURL)
}
