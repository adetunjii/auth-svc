package grpcHandler

import (
	"context"

	"fmt"
	"strconv"

	"github.com/google/uuid"
	"gitlab.com/dh-backend/auth-service/internal/model"
	"gitlab.com/dh-backend/auth-service/internal/util"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ResetPassword(ctx context.Context, request *proto.ResetPasswordRequest) (*proto.ResetPasswordResponse, error) {
	email := request.GetEmail()

	if err := model.IsValidEmail(email); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email format", err)
	}

	user, err := s.Repository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid credentials", err)
	}

	// generate otp
	randomOtp := strconv.Itoa(util.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error("failed to create requestId for reset password", err)
		return nil, status.Errorf(codes.Internal, "failed to create requestId for reset password")
	}

	fmt.Println(requestId.String())
	ev := &model.OtpVerification{
		Otp:   randomOtp,
		Email: user.Email,
	}

	// store otp in cache for 10 minutes using requestId as the key
	if err := s.Redis.SaveOTP(requestId.String(), "RESET_PASSWORD", ev); err != nil {
		s.logger.Error("failed to save otp to redis: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	// create queue message and send to the notification queue
	notificationMessage := model.Notification{
		Otp:              randomOtp,
		User:             *user,
		MessageType:      "reset_password",
		NotificationType: "email",
	}

	s.RabbitMQ.Publish("notification_queue", notificationMessage)

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

	if otpType != proto.OtpType_RESET_PASSWORD {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid OTP type")
	}

	data, err := s.Redis.GetOTP(requestId, model.OtpType(otpType.String()))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "otp has expired, please try again!")
	}

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "otp is incorrect")
	}

	if email != data.Email {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email")
	}

	if err := model.IsValidEmail(email); err == nil {

		_, err = s.Repository.FindUserByEmail(ctx, email)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "invalid credentials", err)
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

	data, err := s.Redis.GetOTP(requestId, "RESET_PASSWORD")
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid requestId")
	}

	if email != data.Email {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email")
	}

	if err := model.IsPasswordValid(newPassword); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password is invalid", err)
	}

	if err := model.IsValidEmail(email); err == nil {

		user, err := s.Repository.FindUserByEmail(ctx, email)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
		}

		userPatch := &model.UserPatch{
			Password: &newPassword,
		}

		user.Patch(userPatch)

		err = s.Repository.UpdateUser(ctx, user.Id, user)
		if err != nil {
			s.logger.Error("couldn't update user's password", err)
			return nil, status.Errorf(codes.Internal, "couldn't update user's password: %v", err)
		}

	}

	return &proto.SetNewPasswordResponse{Message: "Password reset complete. Please proceed to login"}, nil
}
