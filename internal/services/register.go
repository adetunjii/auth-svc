package services

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"fmt"
	"github.com/Adetunjii/protobuf-mono/go/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	user := &proto.User{
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Email:     request.GetEmail(),
		Phone:     request.GetPhoneNumber(),
		PhoneCode: request.GetPhoneCode(),
		Password:  request.GetPassword(),
		Address:   request.GetAddress(),
		State:     request.GetState(),
		Country:   request.GetCountry(),
	}

	userRequest := &proto.CreateUserRequest{User: user}

	_, err := s.UserService.CreateUser(ctx, userRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot register user: %v", err))
		return nil, err
	}

	// generate otp
	// TODO: generate random otps
	randomOtp := "123456"
	requestId, err := uuid.NewRandom()
	if err != nil {
		helpers.LogEvent("ERROR", "failed to create requestId for email verification")
		return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	}

	ev := &models.OtpVerification{
		Otp:   randomOtp,
		Email: user.Email,
	}

	fmt.Println(requestId.String())
	// store otp in cache for 10 minutes using requestId as the key
	if err := s.RedisCache.SaveOTP(requestId.String(), "REG", ev); err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	registerResponse := &proto.RegisterResponse{
		Message:   "An OTP has been sent to your email.",
		RequestId: requestId.String(),
	}

	return registerResponse, nil
}
