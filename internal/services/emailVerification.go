package services

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"dh-backend-auth-sv/internal/proto"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) InitEmailVerification(ctx context.Context, request *proto.InitEmailVerificationRequest) (*proto.InitEmailVerificationResponse, error) {
	email := request.GetEmail()
	otpType := request.GetType()

	getUserRequest := &proto.GetUserDetailsByEmailRequest{Email: email}
	res, err := s.UserService.GetUserDetailsByEmail(ctx, getUserRequest)
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = json.Unmarshal(res.GetResponse(), user)
	if err != nil {
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

	ev := &models.EmailVerification{
		Otp:   randomOtp,
		Email: user.Email,
	}

	fmt.Println(requestId.String())
	// store otp in cache for 10 minutes using requestId as the key
	if err := s.RedisCache.SaveOTP(requestId.String(), otpType, ev); err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	response := &proto.InitEmailVerificationResponse{Message: "An OTP has been sent to your email", RequestId: requestId.String()}
	return response, nil
}

func (s *Server) VerifyEmail(ctx context.Context, request *proto.EmailVerificationRequest) (*proto.EmailVerificationResponse, error) {
	email := request.GetEmail()
	otp := request.GetOtp()
	requestId := request.GetRequestID()
	otpType := request.GetType()

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

	response := &proto.EmailVerificationResponse{Message: "successfully verified email"}
	return response, nil
}
