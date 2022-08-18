package services

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
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

	if user.IsEmailVerified {
		return nil, status.Error(codes.PermissionDenied, "Email already verified")
	}

	// generate otp

	randomOtp := strconv.Itoa(helpers.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		helpers.LogEvent("ERROR", "failed to create requestId for email verification")
		return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	}

	ev := &models.OtpVerification{
		Otp:   randomOtp,
		Email: user.Email,
	}

	// store otp in cache for 10 minutes using requestId as the key
	if err := s.RedisCache.SaveOTP(requestId.String(), models.OtpType(otpType.String()), ev); err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	// create queue message and send to the notification queue
	queueMessage := models.QueueMessage{
		Otp:              randomOtp,
		User:             *user,
		MessageType:      "reg_email_verification",
		NotificationType: "email",
	}

	s.RabbitMQ.Publish("notification_queue", queueMessage)

	response := &proto.InitEmailVerificationResponse{Message: "An OTP has been sent to your email", RequestId: requestId.String()}
	return response, nil
}

func (s *Server) VerifyEmail(ctx context.Context, request *proto.EmailVerificationRequest) (*proto.EmailVerificationResponse, error) {
	email := request.GetEmail()
	otp := request.GetOtp()
	requestId := request.GetRequestID()
	otpType := request.GetType()

	data, err := s.RedisCache.GetOTP(requestId, models.OtpType(otpType.String()))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "otp has expired, please try again!")
	}

	fmt.Println(data)

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "otp is incorrect")
	}

	if email != data.Email {
		helpers.LogEvent("ERROR", fmt.Sprintf("email %s doesn't match %s", email, data.Email))
		return nil, status.Errorf(codes.InvalidArgument, "verification failed")
	}

	user := &models.User{}
	userRequest := proto.GetUserDetailsByEmailRequest{Email: email}

	userDetailsResponse, err := s.UserService.GetUserDetailsByEmail(ctx, &userRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("user not found: %v", err))
		return nil, err
	}

	err = json.Unmarshal(userDetailsResponse.GetResponse(), user)
	if err != nil {
		fmt.Println(err)
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
		return nil, status.Errorf(codes.Internal, "cannot process user info")
	}

	updateUserInfo := proto.UpdateUserInformation{IsEmailVerified: true}
	updateUserRequest := proto.UpdateUserInformationRequest{Id: user.ID, PersonalInformation: &updateUserInfo}

	updateUser, err := s.UserService.UpdateUserInformation(ctx, &updateUserRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot update user %v", updateUser))
		return nil, status.Errorf(codes.Internal, "cannot update user %v", err)
	}

	response := &proto.EmailVerificationResponse{Message: "successfully verified email"}
	return response, nil
}
