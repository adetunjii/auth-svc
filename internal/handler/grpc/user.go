package grpc

import (
	"context"

	"github.com/adetunjii/auth-svc/internal/model"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUserInformation(ctx context.Context, request *proto.UpdateUserInformationRequest) (*proto.UpdateUserInformationResponse, error) {
	userId := request.GetId()
	firstName := request.GetFirstName()
	lastName := request.GetLastName()
	address := request.GetAddress()
	state := request.GetState()

	userPatch := &model.UserPatch{
		FirstName: &firstName,
		LastName:  &lastName,
		Address:   &address,
		State:     &state,
	}

	user, err := s.store.User().FindById(ctx, userId)
	if err != nil {
		s.logger.Error("user not found with err: ", err)
		return nil, status.Errorf(codes.InvalidArgument, "user not found")
	}

	if err := user.Patch(userPatch); err != nil {
		s.logger.Error("failed to update user with err: ", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	if err := s.store.User().Update(ctx, userId, user); err != nil {
		s.logger.Error("failed to update user with err: ", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	// fetch updated user
	updatedUser, err := s.store.User().FindById(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user not found")
	}
	protoUserResponse := &proto.User{
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
		Address:   updatedUser.Address,
	}

	response := &proto.UpdateUserInformationResponse{
		Message: "successfully updated user",
		User:    protoUserResponse,
	}
	return response, nil
}
