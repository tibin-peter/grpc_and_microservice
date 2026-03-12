package main

import (
	"log"
	"net"

	userpb "grpc_and_microservice/proto/user"
	"grpc_and_microservice/user-service/db"
	"grpc_and_microservice/user-service/internal/handler"
	"grpc_and_microservice/user-service/internal/repository"
	"grpc_and_microservice/user-service/internal/service"

	"google.golang.org/grpc"
)

func main() {

	// database
	pg := db.ConnectPostgres()

	// redis
	rdb := db.ConnectRedis()

	// repository
	userRepo := repository.NewRepository(pg)

	// service
	userService := service.NewUserService(userRepo, rdb)

	// handler
	userHandler := handler.NewUserHandler(userService)

	lis, err := net.Listen("tcp",":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	userpb.RegisterUserServiceServer(server, userHandler)

	log.Println("User service running on port 50051")

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}