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

func (s *Server) ResetPassword(ctx context.Context, request *proto.ResetPasswordRequest) (*proto.ResetPasswordResponse, error) {
	email := request.GetEmail()

	getUserByEmailRequest := &proto.GetUserDetailsByEmailRequest{
		Email: email,
	}

	isEmailValid := helpers.IsEmailValid(email)
	if !isEmailValid {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrInvalidEmail, email))
		return nil, status.Errorf(codes.InvalidArgument, "Email is not valid")
	}

	userDetails, err := s.UserService.GetUserDetailsByEmail(ctx, getUserByEmailRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("user not found: %s", err))
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	user := &models.User{}

	err = json.Unmarshal(userDetails.GetResponse(), user)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user: %v", err))
		return nil, status.Errorf(codes.Internal, "cannot unmarshal user")
	}

	// generate otp
	randomOtp := strconv.Itoa(helpers.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		helpers.LogEvent("ERROR", "failed to create requestId for reset password")
		return nil, status.Errorf(codes.Internal, "failed to create requestId for reset password")
	}

	fmt.Println(requestId)

	fmt.Println(requestId.String())
	ev := &models.OtpVerification{
		Otp:   randomOtp,
		Email: user.Email,
	}

	// store otp in cache for 10 minutes using requestId as the key
	if err := s.RedisCache.SaveOTP(requestId.String(), "RESET_PASSWORD", ev); err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	// create queue message and send to the notification queue
	queueMessage := models.QueueMessage{
		Otp:              randomOtp,
		User:             *user,
		MessageType:      "reset_password",
		NotificationType: "email",
	}

	s.RabbitMQ.Publish("notification_queue", queueMessage)

	return &proto.ResetPasswordResponse{
		Message:   "A Password Reset OTP has been sent to your email",
		RequestId: requestId.String(),
	}, nil
}

func (s *Server) VerifyPasswordReset(ctx context.Context, request *proto.VerifyPasswordResetRequest) (*proto.VerifyPasswordResetResponse, error) {
	email := request.GetEmail()
	otp := request.GetOtp()
	requestId := request.GetRequestId()
	otpType := request.GetType()

	user := &models.User{}

	if otpType != proto.OtpType_RESET_PASSWORD {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid OTP type")
	}

	data, err := s.RedisCache.GetOTP(requestId, models.OtpType(otpType.String()))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "otp has expired, please try again!")
	}

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "otp is incorrect")
	}

	if email != data.Email {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email")
	}

	if helpers.IsEmailValid(email) {

		userRequest := proto.GetUserDetailsByEmailRequest{
			Email: email,
		}

		res, err := s.UserService.GetUserDetailsByEmail(context.Background(), &userRequest)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist: %v", err))
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
		}

		err = json.Unmarshal(res.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}
	}

	verifyPasswordResetResponse := &proto.VerifyPasswordResetResponse{
		Message: "Otp verified successfully",
	}
	return verifyPasswordResetResponse, nil
}

func (s *Server) SetNewPassword(ctx context.Context, request *proto.SetNewPasswordRequest) (*proto.SetNewPasswordResponse, error) {
	requestId := request.GetRequestId()
	email := request.GetEmail()
	newPassword := request.GetNewPassword()

	data, err := s.RedisCache.GetOTP(requestId, "RESET_PASSWORD")
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid requestId")
	}

	if email != data.Email {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email")
	}

	sevenOrMore, number, upper, special := helpers.VerifyPassword(newPassword)
	if !sevenOrMore {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrPassword, newPassword))
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

	if len(newPassword) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters long")
	}

	hashedPassword, err := helpers.GenerateHashPassword(newPassword)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrGenerateHashPassword, err.Error()))
	}

	user := &models.User{}

	if helpers.IsEmailValid(email) {
		userRequest := proto.GetUserDetailsByEmailRequest{
			Email: email,
		}

		res, err := s.UserService.GetUserDetailsByEmail(context.Background(), &userRequest)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist: %v", err))
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
		}

		err = json.Unmarshal(res.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}

		updateInfo := &proto.UpdateUserInformationRequest{
			Id: user.ID,
			PersonalInformation: &proto.UpdateUserInformation{
				Password: string(hashedPassword),
			},
		}

		_, err = s.UserService.UpdateUserInformation(ctx, updateInfo)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("couldn't update user's password: %v", err))
			return nil, status.Errorf(codes.Internal, "couldn't update user's password: %v", err)
		}

	}

	return &proto.SetNewPasswordResponse{Message: "Password reset complete. Please proceed to login"}, nil
}
