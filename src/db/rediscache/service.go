package rediscache

import (
	"dh-backend-auth-sv/src/models"
	"encoding/json"
	"golang.org/x/net/context"
	"log"
	"time"
)

func (r *RedisCache) SaveSubChannel(key string, channel *models.User) error {
	client := r.getClient()
	json, err := json.Marshal(channel)
	if err != nil {
		log.Panic(err)
	}

	err = client.Set(context.Background(), key, string(json), r.expires*time.Minute*5).Err()
	if err != nil {
		return err
	}
	_, err = client.Ping(context.Background()).Result()
	return err
}

func (r *RedisCache) GetSubChannel(key string) *models.User {
	client := r.getClient()
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("error getting channel from cache: %v", err)
	}
	var user = &models.User{}
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Println("channel retrieved from cache")
	return user
}

func (r *RedisCache) SaveRoleChannel(key string, channel []byte) error {
	client := r.getClient()
	json, err := json.Marshal(channel)
	if err != nil {
		log.Panic(err)
	}

	err = client.Set(context.Background(), key, string(json), r.expires*time.Minute*5).Err()
	if err != nil {
		return err
	}
	_, err = client.Ping(context.Background()).Result()
	return err
}

func (r *RedisCache) GetRoleChannels(key string) []*models.UserRole {
	client := r.getClient()
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("error getting channel from role cache: %v", err)
	}
	var user []*models.UserRole
	err = json.Unmarshal([]byte(val), user)
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Println("channel retrieved from cache")
	log.Println(user)
	return user
}
