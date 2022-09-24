package grpc

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"gitlab.com/dh-backend/auth-service/internal/model"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	phoneNumber := model.TrimPhoneNumber(request.GetPhoneNumber(), request.GetPhoneCode())

	u := &model.User{
		FirstName:   request.GetFirstName(),
		LastName:    request.GetLastName(),
		Email:       request.GetEmail(),
		PhoneNumber: phoneNumber,
		PhoneCode:   request.GetPhoneCode(),
		Password:    request.GetPassword(),
		Address:     request.GetAddress(),
		State:       request.GetState(),
		Country:     request.GetCountry(),
	}

	oauthId := request.GetOauthId()

	if oauthId != "" {
		// fetch user's email from cache and compare to the one sent
		email, err := s.Redis.GetNewOauthuser(oauthId)
		if err != nil {
			s.logger.Error("user registration failed with err: ", err)
			return nil, status.Errorf(codes.Unauthenticated, "registration failed, Invalid oauth id")
		}

		if email != request.GetEmail() {
			s.logger.Error("user registration failed with err: ", errors.New("emails do not match"))
			return nil, status.Errorf(codes.Unauthenticated, "registration failed. Oauth user not found")
		}
		u.IsEmailVerified = true

	}

	err := s.store.User().Save(ctx, u)
	if err != nil {
		if dbErr := err.(*pgconn.PgError); dbErr != nil && dbErr.Code == "23505" {
			return nil, status.Errorf(codes.InvalidArgument, "user with email / phone number already exists")
		}
		return nil, status.Errorf(codes.InvalidArgument, "invalid parameters")
	}

	registerResponse := &proto.RegisterResponse{
		Message:         "Successfully registered user",
		IsEmailVerified: convertToBooleanString(u.IsEmailVerified),
		IsPhoneVerified: convertToBooleanString(u.IsPhoneVerified),
	}

	return registerResponse, nil
}
