package grpc

import (
	"context"
	"strconv"

	"github.com/adetunjii/auth-svc/internal/model"
	"github.com/adetunjii/auth-svc/internal/util"
	"github.com/google/uuid"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) InitPhoneVerification(ctx context.Context, request *proto.InitPhoneVerificationRequest) (*proto.InitPhoneVerificationResponse, error) {
	phoneCode := request.GetPhoneCode()
	phone := model.TrimPhoneNumber(request.GetPhoneNumber(), phoneCode)
	otpType := request.GetType()

	if otpType != proto.OtpType_REG {
		return nil, status.Errorf(codes.InvalidArgument, "invalid otp type")
	}

	user, err := s.store.User().FindByPhoneNumber(ctx, phone, phoneCode)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials", err)
	}

	if user.IsPhoneVerified {
		return nil, status.Error(codes.PermissionDenied, "Phone number already verified")
	}

	randomOtp := strconv.Itoa(util.RandomOtp())

	requestId, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error("failed to create requestId for Phone Verification", err)
		return nil, status.Errorf(codes.Internal, "failed to create request Id for phone verification")
	}

	ov := &model.OtpVerification{
		Otp:       randomOtp,
		Phone:     user.PhoneNumber,
		PhoneCode: user.PhoneCode,
	}

	if err := s.Redis.SaveOTP(requestId.String(), model.OtpType("REG"), ov); err != nil {
		s.logger.Error("failed to save otp to redis: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	// create queue message and send to the notification queue
	notificationMessage := model.Notification{
		Otp:              randomOtp,
		User:             *user,
		MessageType:      "reg_phone_verification",
		NotificationType: "sms",
	}

	s.RabbitMQ.Publish("notification_queue", notificationMessage)
	response := &proto.InitPhoneVerificationResponse{Message: "An OTP has been sent to your phone number", RequestId: requestId.String()}
	return response, nil
}

func (s *Server) VerifyPhone(ctx context.Context, request *proto.PhoneVerificationRequest) (*proto.PhoneVerificationResponse, error) {
	phoneCode := request.GetPhoneCode()
	phone := model.TrimPhoneNumber(request.GetPhoneNumber(), phoneCode)

	otp := request.GetOtp()
	requestId := request.GetRequestID()
	otpType := request.GetType()

	data, err := s.Redis.GetOTP(requestId, model.OtpType(otpType.String()))
	if err != nil {
		s.logger.Error("invalid requestId", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid request id!")
	}

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "verification failed. Please try again")
	}

	if phone != data.Phone || phoneCode != data.PhoneCode {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid phone number/phone code")
	}

	user, err := s.store.User().FindByPhoneNumber(ctx, phone, phoneCode)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "verification failed!! Invalid Email")
	}

	isPhoneVerified := true
	isActive := true

	updates := &model.User{}
	userPatch := &model.UserPatch{
		IsPhoneVerified: &isPhoneVerified,
		IsActive:        &isActive,
	}

	if err := updates.Patch(userPatch); err != nil {
		s.logger.Error("failed to update user", err)
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	err = s.store.User().Update(ctx, user.Id, updates)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user", err)
	}

	response := &proto.PhoneVerificationResponse{Message: "successfully verified phone number"}
	return response, nil
}
