package main

import (
	"log"
	"net"

	pb "digital-mall/pkg/proto/auth"

	"google.golang.org/grpc"
)

func main() {
	config := LoadConfig()
	db := ConnectDB(config.DatabaseDSN)

	listener, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", config.Port, err)
	}

	grpcServer := grpc.NewServer()
	authService := &AuthService{
		db:     db.Collection("users"),
		config: config,
	}

	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Printf("Auth service is running on port %s", config.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
