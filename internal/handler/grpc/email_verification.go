package grpc

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"gitlab.com/dh-backend/auth-service/internal/model"
	"gitlab.com/dh-backend/auth-service/internal/util"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) InitEmailVerification(ctx context.Context, request *proto.InitEmailVerificationRequest) (*proto.InitEmailVerificationResponse, error) {
	email := request.GetEmail()
	otpType := request.GetType()

	user, err := s.store.User().FindByEmail(ctx, email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user does not exist", err)
	}

	if user.IsEmailVerified {
		return nil, status.Error(codes.PermissionDenied, "email already verified")
	}

	// generate otp
	randomOtp := strconv.Itoa(util.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error("failed to create requestId for email verification", err)
		return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	}

	ev := &model.OtpVerification{
		Otp:   randomOtp,
		Email: user.Email,
	}

	// store otp in cache for 10 minutes using requestId as the key
	if err := s.Redis.SaveOTP(requestId.String(), model.OtpType(otpType.String()), ev); err != nil {
		s.logger.Error("failed to save otp to redis: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	// create queue message and send to the notification queue
	notificationMessage := model.Notification{
		Otp:              randomOtp,
		User:             *user,
		MessageType:      "reg_email_verification",
		NotificationType: "email",
	}

	err = s.RabbitMQ.Publish("notification_queue", notificationMessage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message to notification queue", err)
	}

	response := &proto.InitEmailVerificationResponse{Message: "An OTP has been sent to your email", RequestId: requestId.String()}
	return response, nil
}

func (s *Server) VerifyEmail(ctx context.Context, request *proto.EmailVerificationRequest) (*proto.EmailVerificationResponse, error) {
	email := request.GetEmail()
	otp := request.GetOtp()
	requestId := request.GetRequestID()
	otpType := request.GetType()

	data, err := s.Redis.GetOTP(requestId, model.OtpType(otpType.String()))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "otp has expired, please try again!")
	}

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "otp is incorrect")
	}

	if email != data.Email {
		return nil, status.Errorf(codes.InvalidArgument, "verification failed!! Invalid Email")
	}

	user, err := s.store.User().FindByEmail(ctx, email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "verification failed!! Invalid Email")
	}

	isEmailVerified := true

	updates := &model.User{}
	userPatch := &model.UserPatch{
		IsEmailVerified: &isEmailVerified,
	}

	if err := updates.Patch(userPatch); err != nil {
		s.logger.Error("failed to update user", err)
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	s.store.User().Update(ctx, user.Id, updates)
	response := &proto.EmailVerificationResponse{Message: "successfully verified email"}
	return response, nil
}
