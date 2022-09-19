package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"gitlab.com/dh-backend/auth-service/config"
	"gitlab.com/dh-backend/auth-service/internal/port"
	"gitlab.com/dh-backend/auth-service/internal/services/oauth"
	"gitlab.com/dh-backend/auth-service/internal/services/rabbitmq"
	"gitlab.com/dh-backend/auth-service/internal/services/redis"
	"gitlab.com/dh-backend/auth-service/internal/util"
	"gitlab.com/grpc-buffer/proto/go/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Repository   port.Repository
	Redis        *redis.Redis
	RabbitMQ     *rabbitmq.Connection
	jwtFactory   *util.JwtFactory
	googleClient *oauth.GoogleClient
	logger       port.AppLogger

	proto.UnimplementedAuthServiceServer
}

func New(service *config.Service, logger port.AppLogger) *Server {
	return &Server{
		Repository:   service.Repository,
		Redis:        service.Redis,
		RabbitMQ:     service.RabbitMQ,
		googleClient: service.GoogleClient,
		logger:       logger,
	}
}

func (s *Server) Start(port string) {

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// register services
	proto.RegisterAuthServiceServer(grpcServer, s)
	reflection.Register(grpcServer) // register reflection service on gRPC services

	// graceful shutdown
	go func() {
		s.logger.Info(fmt.Sprintf("server is running on port %v...", port))

		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the services with
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until signal is received
	<-ch
	fmt.Println("stopping the services")

	grpcServer.Stop()

	fmt.Println("closing the listener")
	err = listen.Close()
	if err != nil {
		return
	}

}
