package services

import (
	"context"
	"dh-backend-auth-sv/internal/helpers"
	"dh-backend-auth-sv/internal/proto"
	"fmt"
)

func (s *Server) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	user := &proto.User{
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Email:     request.GetEmail(),
		Phone:     request.GetPhoneNumber(),
		Password:  request.GetPassword(),
		Country:   request.GetCountry(),
	}
	res, err := s.UserService.CreateUser(ctx, &proto.CreateUserRequest{User: user})
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("cannot register user"))
		return nil, err
	}

	registerResponse := &proto.RegisterResponse{
		Message:       res.GetMessage(),
		UserReference: res.GetUserReference(),
	}

	return registerResponse, nil
}
