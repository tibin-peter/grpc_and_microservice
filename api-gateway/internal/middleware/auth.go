package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {

	token := c.Cookies("access_token")

	if token == "" {
		return c.Status(401).JSON("unauthorized")
	}

	return c.Next()
}