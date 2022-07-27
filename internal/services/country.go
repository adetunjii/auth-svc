package services

import (
	"context"
	"dh-backend-auth-sv/internal/proto"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *Server) GetAllCountries(ctx context.Context, req *emptypb.Empty) (*proto.GetAllCountryResponse, error) {
	countries, err := s.DB.GetAllCountries()
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.NotFound, "cannot fetch countries")
	}

	j, err := json.Marshal(countries)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch countries")
	}
	fmt.Println(countries)
	var c []*proto.Country

	err = json.Unmarshal(j, &c)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "cannot fetch countries")
	}

	fmt.Println(c)

	res := &proto.GetAllCountryResponse{Countries: c}

	return res, nil

}
