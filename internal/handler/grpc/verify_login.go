package grpcHandler

import (
	"context"

	//"dh-backend-auth-sv/internal/proto"

	"time"

	"gitlab.com/dh-backend/auth-service/internal/model"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

		user, err = s.Repository.FindUserByEmail(ctx, login)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
		}

	} else {
		phoneCode := req.GetPhoneCode()

		if login != data.Phone && phoneCode != data.PhoneCode {
			return nil, status.Errorf(codes.InvalidArgument, "invalid phone number")
		}

		phone := model.TrimPhoneNumber(login, phoneCode)

		user, err = s.Repository.FindUserByPhoneNumber(ctx, phone, phoneCode)
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
	}

	token, err := s.jwtFactory.CreateToken(ui, exp)
	if err != nil {
		s.logger.Error("failed to create token", err)
		return nil, status.Errorf(codes.Internal, "failed to create token", err)
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
		Token:    token,
		UserInfo: userInfo,
	}

	return loginResponse, nil
}
