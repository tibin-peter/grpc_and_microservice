package handler

import (
	"context"
	"grpc_and_microservice/api-gateway/internal/client"
	"grpc_and_microservice/api-gateway/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	client *client.UserClient
}

func NewUserHandler(client *client.UserClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {

	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON("invalid request")
	}

	_, err := h.client.Register(context.Background(), req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(dto.MessageResponse{Message: "user registered"})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {

	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON("invalid request")
	}

	res, err := h.client.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HTTPOnly: true,
	})

	return c.JSON(dto.MessageResponse{Message: "login successful"})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {

	token := c.Cookies("refresh_token")

	res, err := h.client.Refresh(context.Background(), token)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		HTTPOnly: true,
	})

	return c.JSON(dto.MessageResponse{Message: "token refreshed"})
}

func (h *UserHandler) ValidateToken(c *fiber.Ctx) error {

	token := c.Cookies("access_token")

	res, err := h.client.Validate(context.Background(), token)
	if err != nil || !res.Valid {
		return c.Status(401).JSON("invalid token")
	}

	return c.JSON(dto.ValidateResponse{
		UserID: uint(res.UserId),
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {

	token := c.Cookies("refresh_token")

	_, err := h.client.Logout(context.Background(), token)
	if err != nil {
		return err
	}

	c.ClearCookie("access_token", "refresh_token")

	return c.JSON(dto.MessageResponse{Message: "logged out"})
}