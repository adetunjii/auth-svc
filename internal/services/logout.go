package services

import (
	"context"
	"log"

	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"strings"
)

func (s *Server) Logout(ctx context.Context, request *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	userID := strings.TrimSpace(request.GetUserId())
	if userID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user id cannot be empty")
	}

	err := s.DB.DeleteActivities(userID)
	if err != nil {
		log.Println(err)
	}
	return &proto.LogoutResponse{Message: "activity deleted successfully"}, nil
}
