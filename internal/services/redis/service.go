package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/adetunjii/auth-svc/internal/model"
)

func (r *Redis) SaveOTP(key string, otpType model.OtpType, value any) error {
	json, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error parsing value")
	}

	var redisKey string

	switch otpType {
	case model.LOGIN:
		redisKey = fmt.Sprintf("LOGIN/%s", key)
	case model.REG:
		redisKey = fmt.Sprintf("REG/%s", key)
	case model.RESET_PASSWORD:
		redisKey = fmt.Sprintf("RESET_PASSWORD/%s", key)
	default:
		return fmt.Errorf("invalid otp type")
	}

	return r.client.Set(context.Background(), redisKey, string(json), 10*time.Minute).Err()
}

func (r *Redis) GetOTP(key string, otpType model.OtpType) (*model.OtpVerification, error) {
	response := &model.OtpVerification{}

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

func (r *Redis) SaveNewOauthUser(oauthId string, email string) error {
	redisKey := "OAUTH_REG/" + oauthId
	return r.client.Set(context.Background(), redisKey, email, 10*time.Minute).Err()
}

func (r *Redis) GetNewOauthuser(oauthId string) (string, error) {
	redisKey := "OAUTH_REG/" + oauthId
	return r.client.Get(context.Background(), redisKey).Result()
}

func (r *Redis) SaveRoles(ctx context.Context, roles map[string]interface{}) error {
	return r.client.Set(context.Background(), "roles", roles, 0).Err()
}

func (r *Redis) SavePermissions(ctx context.Context, permissions map[string]interface{}) error {
	return r.client.Set(context.Background(), "permissions", permissions, 0).Err()
}

func (r *Redis) SaveUserRoles(ctx context.Context, userRoles map[string]interface{}) error {
	return r.client.Set(context.Background(), "userRoles", userRoles, 0).Err()
}
