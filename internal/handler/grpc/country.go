package grpc

import (
	"context"
	"encoding/json"

	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAllCountries(ctx context.Context, req *proto.GetAllCountryRequest) (*proto.GetAllCountryResponse, error) {
	countries, err := s.Repository.ListCountries(ctx)
	if err != nil {
		s.logger.Error("cannot fetch countries", err)
		return nil, status.Errorf(codes.NotFound, "cannot fetch countries")
	}

	j, err := json.Marshal(countries)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch countries")
	}

	// unmarshalling into the proto's country struct
	c := []*proto.Country{}

	err = json.Unmarshal(j, &c)
	if err != nil {
		s.logger.Error("cannot unmarshal into country struct", err)
		return nil, status.Errorf(codes.Internal, "cannot fetch countries")
	}

	res := &proto.GetAllCountryResponse{Countries: c}

	return res, nil

}
