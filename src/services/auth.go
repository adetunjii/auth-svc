package services

import (
	"dh-backend-auth-sv/rabbitMQ"
	"dh-backend-auth-sv/src/auth"
	"dh-backend-auth-sv/src/helpers"
	"dh-backend-auth-sv/src/ports"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type Server struct {
	DB         ports.DB
	RedisCache ports.RedisCache
	auth.UnimplementedAuthServiceServer
}

func (s *Server) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	email := strings.TrimSpace(request.GetEmail())
	if !helpers.IsEmailValid(email) {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrInvalidEmail, email))
		return nil, status.Errorf(codes.InvalidArgument, "Email is not valid")
	}
	password := strings.TrimSpace(request.GetPassword())
	sevenOrMore, number, upper, special := helpers.VerifyPassword(password)
	if !sevenOrMore {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrPassword, password))
		return nil, status.Errorf(codes.InvalidArgument, helpers.ErrPassword)
	}
	if !number {
		helpers.LogEvent("ERROR", "Password must contain at least one number")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one number")
	}
	if !upper {
		helpers.LogEvent("ERROR", "Password must contain at least one uppercase letter")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one uppercase letter")
	}
	if !special {
		helpers.LogEvent("ERROR", "Password must contain at least one special character")
		return nil, status.Errorf(codes.InvalidArgument, "Password must contain at least one special character")
	}

	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters long")
	}
	hashedPassword, err := helpers.GenerateHashPassword(password)
	if err != nil {
		helpers.LogEvent("ERROR", fmt.Sprintf("%s: %s", helpers.ErrGenerateHashPassword, err.Error()))
	}

	err = rabbitMQ.PublishToLoginQueue(hashedPassword, email)
	if err != nil {
		return nil, err
	}

	user := s.RedisCache.GetSubChannel(email)
	if user.Email == "" {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if !helpers.CheckPasswordHash(password, []byte(user.HashedPassword)) {
		return nil, status.Error(codes.NotFound, "user password incorrect")
	}

	return &auth.LoginResponse{
		Token: "hello",
	}, nil
}
