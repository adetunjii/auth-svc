package services

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"dh-backend-auth-sv/internal/proto"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) VerifyLogin(ctx context.Context, req *proto.VerifyLoginRequest) (*proto.VerifyLoginResponse, error) {
	otp := req.GetOtp()
	requestId := req.GetRequestId()
	otpType := req.GetType()
	email := req.GetEmail()

	data, err := s.RedisCache.GetOTP(requestId, otpType.String())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "otp has expired, please try again!")
	}

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "otp is incorrect")
	}

	if email != data.Email {
		return nil, status.Errorf(codes.InvalidArgument, "email does not match")
	}

	userRequest := proto.GetUserDetailsByEmailRequest{
		Email: email,
	}
	res, err := s.UserService.GetUserDetailsByEmail(context.Background(), &userRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist!"))
		return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
	}

	user := &models.User{}

	err = json.Unmarshal(res.GetResponse(), user)
	if err != nil {
		fmt.Println(err)
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
		return nil, status.Errorf(codes.Internal, "cannot process user info")
	}

	now := time.Now()
	exp := now.Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": map[string]string{
			"userId": user.ID,
			"roleId": user.RoleID,
		},
		"aud": "proto-service",
		"iss": "proto-service",
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	})
	tokenStr, err := token.SignedString([]byte(email))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	activities := &models.Activities{
		ID:     uuid.New().String(),
		UserID: user.ID,
		Token:  tokenStr,
		Time:   time.Now(),
		Device: string(rune(os.Getpid())),
	}

	err = s.DB.SaveActivities(activities)
	if err != nil {
		log.Printf("err %s", err)
	}

	loginResponse := &proto.VerifyLoginResponse{
		Token: tokenStr,
	}

	return loginResponse, nil
}
