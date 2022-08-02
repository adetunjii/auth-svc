package services

import (
	"context"
	"dh-backend-auth-sv/internal/proto"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *Server) GetAllCountries(ctx context.Context, req *proto.GetAllCountryRequest) (*proto.GetAllCountryResponse, error) {
	countries, err := s.DB.GetAllCountries()
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.NotFound, "cannot fetch countries")
	}

	j, err := json.Marshal(countries)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch countries")
	}
	var c []*proto.Country

	err = json.Unmarshal(j, &c)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "cannot fetch countries")
	}

	res := &proto.GetAllCountryResponse{Countries: c}

	return res, nil

}
