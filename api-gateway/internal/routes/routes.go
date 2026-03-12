package routes

import (
	"grpc_and_microservice/api-gateway/internal/handler"
	"grpc_and_microservice/api-gateway/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {

	api := app.Group("/api")

	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Post("/logout", userHandler.Logout)

	api.Get("/validate", middleware.AuthMiddleware, userHandler.ValidateToken)
}