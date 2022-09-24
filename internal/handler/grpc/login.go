package grpc

import (
	"context"
	"errors"
	"strconv"
	"time"

	"fmt"
	"strings"

	"gitlab.com/dh-backend/auth-service/internal/model"
	"gitlab.com/dh-backend/auth-service/internal/util"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidPhoneCode     = errors.New("invalid phone code")
	ErrNotificationNotSent  = errors.New("verification otp not sent")
	ErrRedisMessageNotSaved = errors.New("message did not save in redis")
)

func (s *Server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {

	login := strings.TrimSpace(request.GetLogin())
	phoneCode := strings.TrimSpace(request.GetPhoneCode())
	password := strings.TrimSpace(request.GetPassword())

	err := model.IsPasswordValid(password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password is invalid", err)
	}

	user, err := s.fetchUserWithEmailOrPhone(login, phoneCode)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := model.ComparePassword(user.Password, password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, ErrInvalidCredentials.Error())
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

	// create and store otp in the cache
	randomOtp := strconv.Itoa(util.RandomOtp())
	requestId, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error("failed to create requestId for email verification", err)
		return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
	}

	var ov *model.OtpVerification

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

	// send message to the notification queue
	notificationMessage := model.Notification{
		Otp:  randomOtp,
		User: *user,
	}

	// publish notification based on login type
	if err := model.IsValidEmail(login); err == nil {

		if rmErr := s.publishNotification(notificationMessage, model.LoginEmailVerification, model.Email); rmErr != nil {
			return nil, status.Errorf(codes.Internal, ErrNotificationNotSent.Error())
		}

		response := &proto.LoginResponse{
			Message:         "An otp has been sent to your email",
			RequestId:       requestId.String(),
			IsEmailVerified: convertToBooleanString(isEmailVerified),
			IsPhoneVerified: convertToBooleanString(isPhoneVerified),
		}

		return response, nil
	} else {

		if rmErr := s.publishNotification(notificationMessage, model.LoginPhoneVerification, model.Sms); rmErr != nil {
			return nil, status.Errorf(codes.Internal, ErrNotificationNotSent.Error())
		}

		response := &proto.LoginResponse{
			Message:         "An otp has been sent to your phone number",
			RequestId:       requestId.String(),
			IsEmailVerified: convertToBooleanString(isEmailVerified),
			IsPhoneVerified: convertToBooleanString(isPhoneVerified),
		}
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
		user, err = s.store.User().FindByEmail(ctx, login)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, ErrInvalidCredentials.Error())
		}
	} else {

		if phoneCode == "" {
			return nil, status.Errorf(codes.InvalidArgument, ErrInvalidPhoneCode.Error())
		}
		phoneNumber := model.TrimPhoneNumber(login, phoneCode)

		user, err = s.store.User().FindByPhoneNumber(ctx, phoneNumber, phoneCode)

		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, ErrInvalidCredentials.Error())
		}
	}

	if err = model.ComparePassword(user.Password, password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, ErrInvalidCredentials.Error())
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
		"country":           user.Country,
	}

	token, err := s.jwtFactory.CreateToken(ui, exp)
	if err != nil {
		s.logger.Error("failed to create token", err)
		return nil, status.Errorf(codes.Internal, "failed to create token")
	}

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
		AccessToken: token,
		UserInfo:    userInfo,
	}
	return loginResponse, nil

}

// OAUTH FLOW
// 1. get access_token from the client
// 2. use access_token to get the user's details from google
// 3. Check if the user already exists in our db
// 4. if exists, create auth token and send back to the client.
// 5. if not, store user details in redis with a requestId,
// 6. send details back to the client, to complete the registration flow
func (s *Server) LoginWithGoogle(ctx context.Context, request *proto.LoginWithGoogleRequest) (*proto.LoginWithGoogleResponse, error) {

	access_token := request.GetGoogleAccessToken()

	userDetails, err := s.googleClient.FetchUserDetails(access_token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user authentication failed")
	}

	if _, ok := userDetails["error"]; ok {
		errorMessage := userDetails["error"].(map[string]interface{})
		s.logger.Error(fmt.Sprintf("failed to authenticate user with err: %s", errorMessage["message"]), nil)
		return nil, status.Errorf(codes.Internal, "user authentication failed")
	}

	email := userDetails["email"].(string)

	user, err := s.store.User().FindByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			newUser := &proto.User{
				Email:     email,
				FirstName: userDetails["given_name"].(string),
				LastName:  userDetails["family_name"].(string),
			}

			oauthId, err := uuid.NewRandom()
			if err != nil {
				s.logger.Error("failed to create requestId for email verification", err)
				return nil, status.Errorf(codes.Internal, "failed to create requestId for email verification")
			}

			if err := s.Redis.SaveNewOauthUser(oauthId.String(), email); err != nil {
				s.logger.Error("failed to register user with err: ", err)
				return nil, status.Errorf(codes.Internal, "cannot register new user")
			}

			resp := &proto.LoginWithGoogleResponse{
				IsNewUser:       "true",
				User:            newUser,
				OauthId:         oauthId.String(),
				IsEmailVerified: convertToBooleanString(userDetails["verified_email"].(bool)),
				IsPhoneVerified: "false",
			}

			return resp, nil
		}
		s.logger.Error("failed to authenticate user with err: ", err)
		return nil, status.Errorf(codes.Unauthenticated, "failed to authenticate user.")
	}

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
		"country":           user.Country,
	}

	token, err := s.jwtFactory.CreateToken(ui, exp)
	if err != nil {
		s.logger.Error("failed to create token", err)
		return nil, status.Errorf(codes.Internal, "failed to create token")
	}

	isVerified := user.IsEmailVerified && user.IsPhoneVerified

	resp := &proto.LoginWithGoogleResponse{
		IsNewUser: "false",
		User: &proto.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.PhoneNumber,
			PhoneCode: user.PhoneCode,
			Email:     user.Email,
		},
		IsEmailVerified: convertToBooleanString(user.IsEmailVerified),
		IsPhoneVerified: convertToBooleanString(user.IsPhoneVerified),
	}

	if isVerified {
		resp.AccessToken = token
		return resp, nil
	}

	return resp, nil

}

func (s *Server) fetchUserWithEmailOrPhone(login string, args ...string) (*model.User, error) {
	user := &model.User{}

	if err := model.IsValidEmail(login); err == nil {
		user, err = s.store.User().FindByEmail(context.Background(), login)
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	} else {

		phoneCode := args[0]

		if phoneCode == "" {
			return nil, ErrInvalidPhoneCode
		}

		phoneNumber := model.TrimPhoneNumber(login, phoneCode)

		user, err = s.store.User().FindByPhoneNumber(context.Background(), phoneNumber, phoneCode)
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	return user, nil
}

func (s *Server) VerifyLogin(ctx context.Context, req *proto.VerifyLoginRequest) (*proto.VerifyLoginResponse, error) {

	otp := req.GetOtp()
	requestId := req.GetRequestId()
	otpType := req.GetType()
	login := req.GetLogin()

	var user *model.User

	data, err := s.Redis.GetOTP(requestId, model.OtpType(otpType.String()))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "otp has expired, please try again!")
	}

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "otp is incorrect")
	}

	if err := model.IsValidEmail(login); err == nil {

		if login != data.Email {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email")
		}

		user, err = s.store.User().FindByEmail(ctx, login)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
		}

	} else {
		phoneCode := req.GetPhoneCode()

		if login != data.Phone && phoneCode != data.PhoneCode {
			return nil, status.Errorf(codes.InvalidArgument, "invalid phone number")
		}

		phone := model.TrimPhoneNumber(login, phoneCode)

		user, err = s.store.User().FindByPhoneNumber(ctx, phone, phoneCode)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "user with this phone number does not exist!")
		}
	}

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
		"country":           user.Country,
	}

	token, err := s.jwtFactory.CreateToken(ui, exp)
	if err != nil {
		s.logger.Error("failed to create token", err)
		return nil, status.Errorf(codes.Internal, "failed to create token")
	}

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
		AccessToken: token,
		UserInfo:    userInfo,
	}

	return loginResponse, nil
}

func (s *Server) publishNotification(notificationMessage model.Notification, messageType model.MessageType, notificationType model.NotificationType) error {
	// create a queue message and push to the notification queue
	notificationMessage.MessageType = messageType
	notificationMessage.NotificationType = notificationType

	return s.RabbitMQ.Publish("notification_queue", notificationMessage)

}

func convertToBooleanString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
