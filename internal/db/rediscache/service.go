package rediscache

import (
	"dh-backend-auth-sv/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
)

func (r *RedisCache) SaveSubChannel(key string, channel *models.User) error {
	client := r.GetClient()
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
	client := r.GetClient()
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
	client := r.GetClient()
	log.Println(string(channel))
	err := client.Set(context.Background(), key, string(channel), r.expires*time.Minute*5).Err()
	if err != nil {
		return err
	}
	_, err = client.Ping(context.Background()).Result()
	return err
}

func (r *RedisCache) GetRoleChannels(key string) []models.UserRole {
	client := r.GetClient()
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("error getting channel from role cache: %v", err)
	}
	var user []models.UserRole
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		log.Printf("error: %v", err)
	}
	log.Println("channel retrieved from cache")
	log.Println(user, "printed user")
	return user
}

func (r *RedisCache) SaveOTP(key string, otpType string, value any) error {
	client := r.GetClient()
	json, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error parsing value")
	}

	redisKey := fmt.Sprintf("%s/%s", otpType, key)

	return client.Set(context.Background(), redisKey, string(json), r.expires*10*time.Minute).Err()
}

func (r *RedisCache) GetOTP(key string, otpType string) (*models.EmailVerification, error) {
	client := r.GetClient()
	//response := &services.EmailVerification{}
	response := &models.EmailVerification{}

	redisKey := fmt.Sprintf("%s/%s", otpType, key)

	val, err := client.Get(context.Background(), redisKey).Result()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(val), response)
	if err != nil {
		return response, err
	}

	return response, nil
}
