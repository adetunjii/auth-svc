package grpc

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
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

	user, err := s.Repository.FindUserById(ctx, userId)
	if err != nil {
		s.logger.Error("user not found with err: ", err)
		return nil, status.Errorf(codes.InvalidArgument, "user not found")
	}

	if err := user.Patch(userPatch); err != nil {
		s.logger.Error("failed to update user with err: ", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	if err := s.Repository.UpdateUser(ctx, userId, user); err != nil {
		s.logger.Error("failed to update user with err: ", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	// fetch updated user
	updatedUser, err := s.Repository.FindUserById(ctx, userId)
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