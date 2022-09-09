package grpcHandler

import (
	"context"

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

	err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parameters", err)
	}

	registerResponse := &proto.RegisterResponse{
		Message: "Successfully registered user",
	}

	return registerResponse, nil
}
