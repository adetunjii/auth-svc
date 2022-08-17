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

func (s *Server) InitPhoneVerification(ctx context.Context, request *proto.InitPhoneVerificationRequest) (*proto.InitPhoneVerificationResponse, error) {
	phoneCode := request.GetPhoneCode()
	phone := helpers.TrimPhoneNumber(request.GetPhoneNumber(), phoneCode)
	otpType := request.GetType()

	if otpType != proto.OtpType_REG {
		return nil, status.Errorf(codes.InvalidArgument, "invalid otp type")
	}

	getUserRequest := &proto.GetUserByPhoneNumberRequest{
		Phone:     phone,
		PhoneCode: phoneCode,
	}

	res, err := s.UserService.GetUserDetailsByPhoneNumber(ctx, getUserRequest)
	if err != nil {
		helpers.LogEvent("ERROR", "user does not exist")
		return nil, err
	}

	user := &models.User{}

	err = json.Unmarshal(res.GetResponse(), user)
	if err != nil {
		helpers.LogEvent("ERROR", "cannot unmarshal user")
		return nil, status.Errorf(codes.Internal, "cannot parse user")
	}

	if user.IsPhoneVerified {
		return nil, status.Error(codes.PermissionDenied, "Phone number already verified")
	}

	randomOtp := strconv.Itoa(helpers.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		helpers.LogEvent("ERROR", "failed to create requestId for Phone Verification")
		return nil, status.Errorf(codes.Internal, "failed to create request Id for phone verification")
	}

	ov := &models.OtpVerification{
		Otp:       randomOtp,
		Phone:     user.PhoneNumber,
		PhoneCode: user.PhoneCode,
	}

	if err := s.RedisCache.SaveOTP(requestId.String(), models.OtpType("REG"), ov); err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("failed to save otp to redis: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	// create queue message and send to the notification queue
	queueMessage := models.QueueMessage{
		Otp:              randomOtp,
		User:             *user,
		MessageType:      "reg_phone_verification",
		NotificationType: "sms",
	}

	s.RabbitMQ.Publish("notification_queue", queueMessage)
	response := &proto.InitPhoneVerificationResponse{Message: "An OTP has been sent to your phone number", RequestId: requestId.String()}
	return response, nil
}

func (s *Server) VerifyPhone(ctx context.Context, request *proto.PhoneVerificationRequest) (*proto.PhoneVerificationResponse, error) {
	phoneCode := request.GetPhoneCode()
	phone := helpers.TrimPhoneNumber(request.GetPhoneNumber(), phoneCode)

	otp := request.GetOtp()
	requestId := request.GetRequestID()
	otpType := request.GetType()

	data, err := s.RedisCache.GetOTP(requestId, models.OtpType(otpType.String()))
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("%v", err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request id!")
	}

	fmt.Println("data from cache >>>>", data)

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "verification failed. Please try again")
	}

	if phone != data.Phone || phoneCode != data.PhoneCode {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid phone number/phone code")
	}

	user := &models.User{}

	userRequest := proto.GetUserByPhoneNumberRequest{
		Phone:     phone,
		PhoneCode: phoneCode,
	}

	userDetailsResponse, err := s.UserService.GetUserDetailsByPhoneNumber(ctx, &userRequest)
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

	updateUserInfo := proto.UpdateUserInformation{IsPhoneVerified: true}
	updateUserRequest := proto.UpdateUserInformationRequest{Id: user.ID, PersonalInformation: &updateUserInfo}

	updateUser, err := s.UserService.UpdateUserInformation(ctx, &updateUserRequest)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot update user %v", updateUser))
		return nil, status.Errorf(codes.Internal, "cannot update user %v", err)
	}

	response := &proto.PhoneVerificationResponse{Message: "successfully verified phone number"}
	return response, nil
}
