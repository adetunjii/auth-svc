package services

import (
	"dh-backend-auth-sv/internal/db/postgres"
	"dh-backend-auth-sv/internal/db/rediscache"
	"dh-backend-auth-sv/internal/proto"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "google.golang.org/grpc/reflection"
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

	redisPass := os.Getenv("REDIS_PASS")
	Addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PORT"))
	rediscache.NewRedisCache(Addr, redisPass, 10, 15*time.Second)
	redisCache := &rediscache.RedisCache{}

	// connect to user service via gRPC
	// TODO: introduce service discovery here
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_URL"), grpc.WithInsecure())
	fmt.Println(os.Getenv("USER_SERVICE_URL"))
	if err != nil {
		log.Printf("cannot connect to user service: %v", err)
	}

	defer conn.Close()
	log.Println("connected to user service....")

	userService := proto.NewUserServiceClient(conn)

	pd := &Server{
		DB:          db,
		RedisCache:  redisCache,
		UserService: userService,
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
	proto.RegisterAuthServiceServer(ser, pd)
	reflection.Register(ser) // register reflection service on gRPC services

	log.Printf("server is running on port %v...", PORT)

	// graceful shutdown
	go func() {
		fmt.Printf("service started on port: %s", PORT)
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
