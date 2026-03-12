package main

import (
	"log"

	"grpc_and_microservice/api-gateway/internal/client"
	"grpc_and_microservice/api-gateway/internal/handler"
	"grpc_and_microservice/api-gateway/internal/routes"
	userpb "grpc_and_microservice/proto/user"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:50051",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	userGrpc := userpb.NewUserServiceClient(conn)

	userClient := client.NewUserClient(userGrpc)

	userHandler := handler.NewUserHandler(userClient)

	app := fiber.New()


	routes.SetupRoutes(app, userHandler)

	log.Println("API Gateway running on port 8080")


	app.Listen(":8081")
}

