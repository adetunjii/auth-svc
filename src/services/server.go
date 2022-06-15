package services

import (
	"dh-backend-auth-sv/src/auth"
	"dh-backend-auth-sv/src/db/postgres"
	"dh-backend-auth-sv/src/db/rediscache"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func Start() {
	// Create a new services
	//helpers.InitializeLogDir()

	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if PORT == ":" {
		PORT += "8080"
	}
	db := &postgres.PostgresDB{}
	db.Init()

	Addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PORT"))
	rediscache.NewRedisCache(Addr, 10, 15*time.Second)
	redisCache := &rediscache.RedisCache{}

	pd := &Server{
		DB:         db,
		RedisCache: redisCache,
	}

	go pd.SubscribeToLoginQueue()
	go pd.SubscribeToRoleQueue()

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%v", PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	ser := grpc.NewServer(opts...)

	// Register the service with the gRPC services

	// register services
	auth.RegisterAuthServiceServer(ser, pd)
	reflection.Register(ser) // register reflection service on gRPC services

	// graceful shutdown
	go func() {
		fmt.Println("starting services ...")
		if err := ser.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the services with
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until signal is received
	<-ch
	fmt.Println("stopping the services")

	ser.Stop()

	fmt.Println("closing the listener")

	fmt.Println("database connection closed")
	err = listen.Close()
	if err != nil {
		return
	}
	fmt.Println("End of program")
}
