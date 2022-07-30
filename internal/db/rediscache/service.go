package rediscache

import (
	"dh-backend-auth-sv/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
)

func (r *Redis) SaveSubChannel(key string, channel *models.User) error {

	json, err := json.Marshal(channel)
	if err != nil {
		log.Panic(err)
	}

	err = r.client.Set(context.Background(), key, string(json), r.expiry*time.Minute*5).Err()
	if err != nil {
		return err
	}
	_, err = r.client.Ping(context.Background()).Result()
	return err
}

func (r *Redis) GetSubChannel(key string) *models.User {
	val, err := r.client.Get(context.Background(), key).Result()
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

func (r *Redis) SaveRoleChannel(key string, channel []byte) error {
	log.Println(string(channel))
	err := r.client.Set(context.Background(), key, string(channel), r.expiry*time.Minute*5).Err()
	if err != nil {
		return err
	}
	_, err = r.client.Ping(context.Background()).Result()
	return err
}

func (r *Redis) GetRoleChannels(key string) []models.UserRole {
	val, err := r.client.Get(context.Background(), key).Result()
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

func (r *Redis) SaveOTP(key string, otpType string, value any) error {
	json, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error parsing value")
	}

	redisKey := fmt.Sprintf("%s/%s", otpType, key)

	return r.client.Set(context.Background(), redisKey, string(json), 10*time.Minute).Err()
}

func (r *Redis) GetOTP(key string, otpType string) (*models.EmailVerification, error) {
	//response := &services.EmailVerification{}
	response := &models.EmailVerification{}

	redisKey := fmt.Sprintf("%s/%s", otpType, key)

	val, err := r.client.Get(context.Background(), redisKey).Result()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(val), response)
	if err != nil {
		return response, err
	}

	return response, nil
}
