package services

import (
	"context"
	"dh-backend-auth-sv/src/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

func (s *Server) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	userID := strings.TrimSpace(request.GetUserId())
	if userID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id cannot be empty")
	}

	err := s.DB.DeleteActivities(userID)
	if err != nil {
		log.Println(err)
	}
	return &auth.LogoutResponse{Message: "activity deleted successfully"}, nil
}
