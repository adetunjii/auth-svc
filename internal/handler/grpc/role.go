package grpc

import (
	"context"
	"encoding/json"

	"github.com/adetunjii/auth-svc/internal/model"
	"github.com/jackc/pgconn"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRole(ctx context.Context, request *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error) {
	role := request.GetRole()

	if role == "" {
		return nil, status.Error(codes.InvalidArgument, "role title cannot be empty")
	}

	arg := &model.Role{
		Title: role,
	}

	if err := s.store.Role().Create(ctx, arg); err != nil {

		if dbErr, ok := err.(*pgconn.PgError); ok && dbErr.Code == "23505" {
			return nil, status.Error(codes.InvalidArgument, "role already exists")
		}

		s.logger.Error("failed to create role with err: ", err)
		return nil, status.Error(codes.Internal, "failed to create role")
	}

	response := &proto.CreateRoleResponse{
		Message: "Role created successfully",
	}

	return response, nil
}

func (s *Server) GetAllRoles(ctx context.Context, request *proto.GetAllRolesRequest) (*proto.GetAllRolesResponse, error) {

	roles, err := s.store.Role().List(ctx, nil, 1, 20)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to return roles")
	}

	bytes, err := json.Marshal(roles)
	if err != nil {
		s.logger.Error("failed to marshal roles to json with err: ", err)
		return nil, status.Errorf(codes.Internal, "failed to marshal roles to json")
	}

	rolesList := []*proto.Role{}
	if err := json.Unmarshal(bytes, &rolesList); err != nil {
		s.logger.Error("failed to unmarshal roles into json with err: ", err)
		return nil, status.Errorf(codes.Internal, "failed to unmarshal roles into json")
	}

	response := &proto.GetAllRolesResponse{
		Roles: rolesList,
	}

	return response, nil
}

// TODO: make sure that the role to be deleted isn't assigned to any user
func (s *Server) DeleteRole(ctx context.Context, request *proto.DeleteRoleRequest) (*proto.DeleteRoleResponse, error) {
	roleId := request.GetRoleId()

	if err := s.store.Role().Delete(ctx, roleId); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete role")
	}

	response := &proto.DeleteRoleResponse{
		Message: "Role deleted successfully",
	}

	return response, nil

}

// func (s *Server) UpdateRole(ctx context.Context, in *proto.UpdateRoleRequest) (*proto.UpdateRoleResponse, error)
// func (s *Server) AssignPermissionsToRole(ctx context.Context, in *proto.AssignPermissionsToRoleRequest) (*proto.AssignPermissionsToRoleResponse, error)
// func (s *Server) RemovePermissionsFromRole(ctx context.Context, in *proto.RemovePermissionsFromRoleRequest) (*proto.RemovePermissionsFromRoleResponse, error)
