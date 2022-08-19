package services

import (
	"context"
	"crypto/rsa"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"dh-backend-auth-sv/internal/ports"
	"dh-backend-auth-sv/internal/rabbitMQ"
	"log"
	"os"
	"strconv"
	"time"

	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/google/uuid"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	DB         ports.DB
	RedisCache ports.RedisCache
	RabbitMQ   *rabbitMQ.RabbitMQ
	proto.UnimplementedAuthServiceServer
	jwtKey      *rsa.PrivateKey
	UserService proto.UserServiceClient
}

func (s *Server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	login := strings.TrimSpace(request.GetLogin())
	phoneCode := strings.TrimSpace(request.GetPhoneCode())

	password := strings.TrimSpace(request.GetPassword())
	sevenOrMore, number, upper, special := helpers.VerifyPassword(password)
	if !sevenOrMore {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrPassword, password))
		return nil, status.Errorf(codes.InvalidArgument, helpers.ErrPassword)
	}
	if !number {
		helpers.LogEvent("ERROR", "Password must contain at least one number")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one number")
	}
	if !upper {
		helpers.LogEvent("ERROR", "Password must contain at least one uppercase letter")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one uppercase letter")
	}
	if !special {
		helpers.LogEvent("ERROR", "Password must contain at least one special character")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one special character")
	}

	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters long")
	}
	_, err := helpers.GenerateHashPassword(password)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrGenerateHashPassword, err.Error()))
	}

	var userByEmailResponse *proto.GetUserDetailsByEmailResponse
	var userByPhoneResponse *proto.GetUserByPhoneNumberResponse
	user := &models.User{}

	if helpers.IsEmailValid(login) {
		userRequestByEmail := proto.GetUserDetailsByEmailRequest{
			Email: login,
		}

		userByEmailResponse, err = s.UserService.GetUserDetailsByEmail(context.Background(), &userRequestByEmail)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist: %v", err))
			return nil, status.Errorf(codes.NotFound, "invalid user")
		}

		err = json.Unmarshal(userByEmailResponse.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}

	} else {

		if phoneCode == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid phone code")
		}

		userRequestByPhone := proto.GetUserByPhoneNumberRequest{
			Phone:     helpers.TrimPhoneNumber(login, phoneCode),
			PhoneCode: phoneCode,
		}
		userByPhoneResponse, err = s.UserService.GetUserDetailsByPhoneNumber(context.Background(), &userRequestByPhone)
		if err != nil {
			helpers.LogEvent("ERROR", "user with this phone number does not exist")
			return nil, status.Errorf(codes.NotFound, "invalid user")
		}

		err = json.Unmarshal(userByPhoneResponse.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}
	}

	if !helpers.CheckPasswordHash(password, []byte(user.HashedPassword)) {
		return nil, status.Error(codes.NotFound, "invalid credentials")
	}

	isEmailVerified := user.IsEmailVerified
	isPhoneVerified := user.IsPhoneVerified

	if !isEmailVerified || !isPhoneVerified {
		helpers.LogEvent("ERROR", "Please verify account to proceed!")

		response := &proto.LoginResponse{
			Message:         "Please verify your email to proceed",
			IsEmailVerified: convertToBooleanString(isEmailVerified),
			IsPhoneVerified: convertToBooleanString(isPhoneVerified),
			User: &proto.User{
				Email:     user.Email,
				Phone:     user.PhoneNumber,
				PhoneCode: user.PhoneCode,
			},
		}
		return response, nil

	}

	randomOtp := strconv.Itoa(helpers.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		helpers.LogEvent("ERROR", "failed to create requestId for email verification")
		return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	}
	fmt.Println(randomOtp)
	ov := &models.OtpVerification{}

	if helpers.IsEmailValid(login) {
		ov = &models.OtpVerification{
			Otp:   randomOtp,
			Email: user.Email,
		}
	} else {
		ov = &models.OtpVerification{
			Otp:       randomOtp,
			Phone:     login,
			PhoneCode: phoneCode,
		}
	}

	// store otp in cache for 10 minutes using requestId as the key
	if err := s.RedisCache.SaveOTP(requestId.String(), "LOGIN", ov); err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	if helpers.IsEmailValid(login) {

		// create a queue message and push to the notification queue
		queueMessage := models.QueueMessage{
			Otp:              randomOtp,
			User:             *user,
			MessageType:      "login_email_verification",
			NotificationType: "email",
		}

		s.RabbitMQ.Publish("notification_queue", queueMessage)

		response := &proto.LoginResponse{
			Message:   "An otp has been sent to your email",
			RequestId: requestId.String(),
		}
		return response, nil
	} else {

		queueMessage := models.QueueMessage{
			Otp:              randomOtp,
			User:             *user,
			MessageType:      "login_phone_verification",
			NotificationType: "sms",
		}

		s.RabbitMQ.Publish("notification_queue", queueMessage)

		response := &proto.LoginResponse{
			Message:         "An otp has been sent to your phone number",
			RequestId:       requestId.String(),
			IsEmailVerified: convertToBooleanString(isEmailVerified),
			IsPhoneVerified: convertToBooleanString(isPhoneVerified),
		}

		return response, nil
	}
}

// LoginNoVerification should only be called after email and phone
// verification are completed on sign up.
func (s *Server) LoginNoVerification(ctx context.Context, request *proto.LoginRequest) (*proto.VerifyLoginResponse, error) {
	login := strings.TrimSpace(request.GetLogin())
	phoneCode := strings.TrimSpace(request.GetPhoneCode())

	password := strings.TrimSpace(request.GetPassword())
	sevenOrMore, number, upper, special := helpers.VerifyPassword(password)
	if !sevenOrMore {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrPassword, password))
		return nil, status.Errorf(codes.InvalidArgument, helpers.ErrPassword)
	}
	if !number {
		helpers.LogEvent("ERROR", "Password must contain at least one number")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one number")
	}
	if !upper {
		helpers.LogEvent("ERROR", "Password must contain at least one uppercase letter")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one uppercase letter")
	}
	if !special {
		helpers.LogEvent("ERROR", "Password must contain at least one special character")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one special character")
	}

	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters long")
	}

	//fmt.Println(hashedPassword)

	var userByEmailResponse *proto.GetUserDetailsByEmailResponse
	var userByPhoneResponse *proto.GetUserByPhoneNumberResponse
	var err error
	user := &models.User{}

	if helpers.IsEmailValid(login) {
		userRequestByEmail := proto.GetUserDetailsByEmailRequest{
			Email: login,
		}

		userByEmailResponse, err = s.UserService.GetUserDetailsByEmail(context.Background(), &userRequestByEmail)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist: %v", err))
			return nil, status.Errorf(codes.NotFound, "invalid user")
		}

		err = json.Unmarshal(userByEmailResponse.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}

	} else {

		if phoneCode == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid phone code")
		}

		userRequestByPhone := proto.GetUserByPhoneNumberRequest{
			Phone:     login,
			PhoneCode: phoneCode,
		}
		userByPhoneResponse, err = s.UserService.GetUserDetailsByPhoneNumber(context.Background(), &userRequestByPhone)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this phone number does not exist"))
			return nil, status.Errorf(codes.NotFound, "invalid user")
		}

		err = json.Unmarshal(userByPhoneResponse.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}

	}

	if !helpers.CheckPasswordHash(password, []byte(user.HashedPassword)) {
		return nil, status.Error(codes.NotFound, "invalid credentials")
	}

	fmt.Println(user)

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

	tokenStr, err := token.SignedString([]byte(login))
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

	userInfo := &proto.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.PhoneNumber,
		PhoneCode: user.PhoneCode,
		Address:   user.Address,
		State:     user.State,
		Country:   user.Country,
	}

	loginResponse := &proto.VerifyLoginResponse{
		Token:    tokenStr,
		UserInfo: userInfo,
	}
	return loginResponse, nil

}

func convertToBooleanString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
