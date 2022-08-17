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

func (r *Redis) SaveOTP(key string, otpType models.OtpType, value any) error {
	json, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error parsing value")
	}

	var redisKey string

	switch otpType {
	case models.LOGIN:
		redisKey = fmt.Sprintf("LOGIN/%s", key)
	case models.REG:
		redisKey = fmt.Sprintf("REG/%s", key)
	case models.RESET_PASSWORD:
		redisKey = fmt.Sprintf("RESET_PASSWORD/%s", key)
	default:
		return fmt.Errorf("invalid otp type")
	}

	return r.client.Set(context.Background(), redisKey, string(json), 10*time.Minute).Err()
}

func (r *Redis) GetOTP(key string, otpType models.OtpType) (*models.OtpVerification, error) {
	response := &models.OtpVerification{}

	redisKey := fmt.Sprintf("%s/%s", otpType, key)

	val, err := r.client.Get(context.Background(), redisKey).Result()
	if err != nil {
		return response, err
	}

	err = json.Unmarshal([]byte(val), response)
	if err != nil {
		return response, err
	}
	// del, err := r.client.Del(context.Background(), redisKey).Result()
	// if err != nil {
	// 	return response, err
	// }
	// fmt.Println(del)

	return response, nil
}
