package services

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	"encoding/json"
	"fmt"

	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	phoneNumber := helpers.TrimPhoneNumber(request.GetPhoneNumber(), request.GetPhoneCode())

	user := &proto.User{
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Email:     request.GetEmail(),
		Phone:     phoneNumber,
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

	userByEmailRequest := &proto.GetUserDetailsByEmailRequest{
		Email: user.Email,
	}

	userByEmailResponse, err := s.UserService.GetUserDetailsByEmail(ctx, userByEmailRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("user not saved: %v", err))
		return nil, status.Errorf(codes.Internal, "user not saved", err)
	}

	res := userByEmailResponse.GetResponse()

	userObject := &models.User{}
	err = json.Unmarshal(res, userObject)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("error unmarshalling user %v", err))
		return nil, status.Errorf(codes.Internal, "error unmarshalling user", err)
	}

	// generate otp
	// randomOtp := strconv.Itoa(helpers.RandomOtp())
	// requestId, err := uuid.NewRandom()
	// if err != nil {
	// 	helpers.LogEvent("ERROR", "failed to create requestId for email verification")
	// 	return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	// }

	// ev := &models.OtpVerification{
	// 	Otp:   randomOtp,
	// 	Email: user.Email,
	// }

	// // store otp in cache for 10 minutes using requestId as the key
	// if err := s.RedisCache.SaveOTP(requestId.String(), "REG", ev); err != nil {
	// 	helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
	// 	return nil, status.Errorf(codes.Internal, "failed to save otp")
	// }

	// create queue message and send to the notification queue
	// queueMessage := models.QueueMessage{
	// 	Otp:              randomOtp,
	// 	User:             *userObject,
	// 	MessageType:      "reg_email_verification",
	// 	NotificationType: "email",
	// }

	// s.RabbitMQ.Publish("notification_queue", queueMessage)

	registerResponse := &proto.RegisterResponse{
		Message: "An OTP has been sent to your email.",
		// RequestId: requestId.String(),
	}

	return registerResponse, nil
}
