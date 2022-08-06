package services

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/models"
	//"dh-backend-auth-sv/internal/proto"
	"encoding/json"
	"fmt"
	"github.com/Adetunjii/protobuf-mono/go/pkg/proto"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"
)

func (s *Server) VerifyLogin(ctx context.Context, req *proto.VerifyLoginRequest) (*proto.VerifyLoginResponse, error) {

	otp := req.GetOtp()
	requestId := req.GetRequestId()
	otpType := req.GetType()
	login := req.GetLogin()

	user := &models.User{}

	data, err := s.RedisCache.GetOTP(requestId, otpType.String())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "otp has expired, please try again!")
	}

	if otp != data.Otp {
		return nil, status.Errorf(codes.InvalidArgument, "otp is incorrect")
	}

	fmt.Println(helpers.IsEmailValid(login))

	if helpers.IsEmailValid(login) {

		if login != data.Email {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email")
		}

		userRequest := proto.GetUserDetailsByEmailRequest{
			Email: login,
		}

		res, err := s.UserService.GetUserDetailsByEmail(context.Background(), &userRequest)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this email does not exist!"))
			return nil, status.Errorf(codes.NotFound, "user with this email does not exist!")
		}

		err = json.Unmarshal(res.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}

	} else {
		phoneCode := req.GetPhoneCode()

		if login != data.Phone && phoneCode != data.PhoneCode {
			return nil, status.Errorf(codes.InvalidArgument, "invalid phone number")
		}

		userRequestByPhone := proto.GetUserByPhoneNumberRequest{
			Phone:     login,
			PhoneCode: phoneCode,
		}

		res, err := s.UserService.GetUserDetailsByPhoneNumber(ctx, &userRequestByPhone)
		if err != nil {
			helpers.LogEvent("ERROR", fmt.Sprintf("user with this phone number does not exist!"))
			return nil, status.Errorf(codes.NotFound, "user with this phone number does not exist!")
		}

		err = json.Unmarshal(res.GetResponse(), user)
		if err != nil {
			fmt.Println(err)
			helpers.LogEvent("ERROR", fmt.Sprintf("cannot unmarshal user %v", err))
			return nil, status.Errorf(codes.Internal, "cannot process user info")
		}
	}

	now := time.Now()
	exp := now.Add(time.Hour * 24)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": map[string]string{
			"userId": user.ID,
			"roleId": user.RoleID,
		},
		"aud": "proto-service",
		"iss": "proto-service",
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	})

	tokenStr, err := token.SignedString([]byte(login))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	activities := &models.Activities{
		ID:     uuid.New().String(),
		UserID: user.ID,
		Token:  tokenStr,
		Time:   time.Now(),
		Device: string(rune(os.Getpid())),
	}

	err = s.DB.SaveActivities(activities)
	if err != nil {
		log.Printf("err %s", err)
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
		Token:    tokenStr,
		UserInfo: userInfo,
	}

	return loginResponse, nil
}

func toUserProtoObject(a interface{}) (*proto.User, error) {
	var user *proto.User

	j, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, user)
	if err != nil {
		return nil, err
	}

	return user, nil

}
