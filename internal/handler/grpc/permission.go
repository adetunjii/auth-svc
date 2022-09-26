package grpc

import (
	"context"
	"encoding/json"

	"github.com/adetunjii/auth-svc/internal/model"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreatePermission(ctx context.Context, req *proto.CreatePermissionRequest) (*proto.CreatePermissionResponse, error) {
	permission := req.GetPermission()
	description := req.GetDescription()

	if permission == "" {
		return nil, status.Errorf(codes.InvalidArgument, "permission cannot be empty")
	}

	perm := &model.Permission{
		Name:        permission,
		Description: description,
	}

	err := s.store.Permission().Create(ctx, perm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create permission")
	}

	resp := &proto.CreatePermissionResponse{
		Message: "Permission created successfully",
	}

	return resp, nil
}

// TODO: proto to include page and size
func (s *Server) GetAllPermissions(ctx context.Context, req *proto.GetAllPermissionsRequest) (*proto.GetAllPermissionsResponse, error) {

	permissions, err := s.store.Permission().List(ctx, nil, 1, 20)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch permissions")
	}

	bytes, err := json.Marshal(permissions)
	if err != nil {
		s.logger.Error("failed to marshal permissions array with err: ", err)
		return nil, status.Errorf(codes.Internal, "cannot fetch permissions")
	}

	permissionsList := []*proto.Permission{}

	if err := json.Unmarshal(bytes, &permissionsList); err != nil {
		s.logger.Error("failed to unmarshal permissions array with err: ", err)
		return nil, status.Errorf(codes.Internal, "cannot fetch permissions")
	}

	resp := &proto.GetAllPermissionsResponse{
		Permissions: permissionsList,
	}

	return resp, nil
}
