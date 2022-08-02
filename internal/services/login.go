package services

import (
	"context"
	"crypto/rsa"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"dh-backend-auth-sv/internal/ports"
	"dh-backend-auth-sv/internal/proto"
	"encoding/json"
	"fmt"
	//"github.com/Adetunjii/protobuf-mono/go/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type Server struct {
	DB         ports.DB
	RedisCache ports.RedisCache
	proto.UnimplementedAuthServiceServer
	jwtKey      *rsa.PrivateKey
	UserService proto.UserServiceClient
}

func (s *Server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	login := strings.TrimSpace(request.GetLogin())

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
	hashedPassword, err := helpers.GenerateHashPassword(password)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrGenerateHashPassword, err.Error()))
	}

	fmt.Println(hashedPassword)

	var userByEmailResponse *proto.GetUserDetailsByEmailResponse
	var userByPhoneResponse *proto.GetUserByPhoneNumberResponse
	user := &models.User{}

	if helpers.IsEmailValid(login) {
		userRequestByEmail := proto.GetUserDetailsByEmailRequest{
			Email: login,
		}

		userByEmailResponse, err = s.UserService.GetUserDetailsByEmail(context.Background(), &userRequestByEmail)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist!"))
			return nil, status.Errorf(codes.NotFound, "invalid user")
		}

		err = json.Unmarshal(userByEmailResponse.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}

	} else {
		userRequestByPhone := proto.GetUserByPhoneNumberRequest{
			Phone: login,
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
		return nil, status.Error(codes.NotFound, "user password incorrect")
	}

	randomOtp := "123456"
	requestId, err := uuid.NewRandom()
	if err != nil {
		helpers.LogEvent("ERROR", "failed to create requestId for email verification")
		return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	}

	ev := &models.EmailVerification{
		Otp:   randomOtp,
		Email: user.Email,
	}

	fmt.Println(requestId.String())
	// store otp in cache for 10 minutes using requestId as the key
	if err := s.RedisCache.SaveOTP(requestId.String(), "LOGIN", ev); err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	response := &proto.LoginResponse{
		Message:   "An otp has been sent to your email",
		RequestId: requestId.String(),
	}
	return response, nil
}
