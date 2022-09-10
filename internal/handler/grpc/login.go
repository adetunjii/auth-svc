package grpcHandler

import (
	"context"
	"strconv"
	"time"

	"fmt"
	"strings"

	"gitlab.com/dh-backend/auth-service/internal/model"
	"gitlab.com/dh-backend/auth-service/internal/util"

	"github.com/google/uuid"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {

	login := strings.TrimSpace(request.GetLogin())
	phoneCode := strings.TrimSpace(request.GetPhoneCode())
	password := strings.TrimSpace(request.GetPassword())

	err := model.IsPasswordValid(password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password is invalid", err)
	}

	user := &model.User{}

	if err := model.IsValidEmail(login); err == nil {
		user, err = s.Repository.FindUserByEmail(ctx, login)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid credentials", err)
		}
	} else {

		if phoneCode == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid phone code")
		}

		phoneNumber := model.TrimPhoneNumber(login, phoneCode)

		user, err = s.Repository.FindUserByPhoneNumber(ctx, phoneNumber, phoneCode)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid credentials", err)
		}

	}

	if err := model.ComparePassword(user.Password, password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials", err)
	}

	isEmailVerified := user.IsEmailVerified
	isPhoneVerified := user.IsPhoneVerified

	if !isEmailVerified || !isPhoneVerified {

		response := &proto.LoginResponse{
			Message:         "Please verify your email to proceed",
			IsEmailVerified: convertToBooleanString(isEmailVerified),
			IsPhoneVerified: convertToBooleanString(isPhoneVerified),
			User: &proto.User{
				Email:     user.Email,
				Phone:     user.PhoneNumber,
				PhoneCode: user.PhoneCode,
			},
		}
		return response, nil
	}

	randomOtp := strconv.Itoa(util.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error("failed to create requestId for email verification", err)
		return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	}

	ov := &model.OtpVerification{}

	if err := model.IsValidEmail(login); err == nil {
		ov = &model.OtpVerification{
			Otp:   randomOtp,
			Email: user.Email,
		}
	} else {
		ov = &model.OtpVerification{
			Otp:       randomOtp,
			Phone:     login,
			PhoneCode: phoneCode,
		}
	}

	// store otp in cache for 10 minutes using requestId as the key
	if err := s.Redis.SaveOTP(requestId.String(), "LOGIN", ov); err != nil {
		s.logger.Error("failed to save otp to redis: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to save otp")
	}

	notificationMessage := model.Notification{
		Otp:  randomOtp,
		User: *user,
	}

	if err := model.IsValidEmail(login); err == nil {

		// create a queue message and push to the notification queue
		notificationMessage.MessageType = "login_email_verification"
		notificationMessage.NotificationType = "email"

		err := s.RabbitMQ.Publish("notification_queue", notificationMessage)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "verification alert not sent")
		}

		response := &proto.LoginResponse{
			Message:         "An otp has been sent to your email",
			RequestId:       requestId.String(),
			IsEmailVerified: convertToBooleanString(isEmailVerified),
			IsPhoneVerified: convertToBooleanString(isPhoneVerified),
		}
		fmt.Println(response)

		return response, nil
	} else {

		notificationMessage.MessageType = "login_phone_verification"
		notificationMessage.NotificationType = "sms"

		err := s.RabbitMQ.Publish("notification_queue", notificationMessage)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "verification alert not sent!")
		}

		response := &proto.LoginResponse{
			Message:         "An otp has been sent to your phone number",
			RequestId:       requestId.String(),
			IsEmailVerified: convertToBooleanString(isEmailVerified),
			IsPhoneVerified: convertToBooleanString(isPhoneVerified),
		}

		fmt.Println(response)
		return response, nil
	}
}

// LoginNoVerification should only be called after email and phone
// verification are completed on sign up.
func (s *Server) LoginNoVerification(ctx context.Context, request *proto.LoginRequest) (*proto.VerifyLoginResponse, error) {
	login := strings.TrimSpace(request.GetLogin())
	phoneCode := strings.TrimSpace(request.GetPhoneCode())
	password := request.GetPassword()

	var err error
	var user *model.User

	if err := model.IsValidEmail(login); err == nil {
		user, err = s.Repository.FindUserByEmail(ctx, login)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "invalid credentials", err)
		}
	} else {

		if phoneCode == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid phone code")
		}
		phoneNumber := model.TrimPhoneNumber(login, phoneCode)

		user, err = s.Repository.FindUserByPhoneNumber(ctx, phoneNumber, phoneCode)

		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid credentials", err)
		}
	}

	if err = model.ComparePassword(user.Password, password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials", err)
	}

	// token expires in 24 hours
	exp := time.Hour * 24

	ui := map[string]interface{}{
		"id":                user.Id,
		"email":             user.Email,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
		"is_email_verified": user.IsEmailVerified,
		"is_active":         user.IsActive,
		"is_phone_verified": user.IsPhoneVerified,
		"role_id":           user.RoleId,
	}

	token, err := s.jwtFactory.CreateToken(ui, exp)
	if err != nil {
		s.logger.Error("failed to create token", err)
		return nil, status.Errorf(codes.Internal, "failed to create token", err)
	}

	// save activities
	// activities := &models.Activities{
	// 	ID:     uuid.New().String(),
	// 	UserID: user.ID,
	// 	Token:  tokenStr,
	// 	Time:   time.Now(),
	// 	Device: string(rune(os.Getpid())),
	// }

	// err = s.DB.SaveActivities(activities)
	// if err != nil {
	// 	log.Printf("err %s", err)
	// }

	userInfo := &proto.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.PhoneNumber,
		PhoneCode: user.PhoneCode,
		Address:   user.Address,
		State:     user.State,
		Country:   user.Country,
	}

	loginResponse := &proto.VerifyLoginResponse{
		Token:    token,
		UserInfo: userInfo,
	}
	return loginResponse, nil

}

func convertToBooleanString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
