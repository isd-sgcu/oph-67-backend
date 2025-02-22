package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/usecase"
	"github.com/isd-sgcu/oph-67-backend/utils"
)

// AuthMiddleware verifies the JWT from the Authorization header
func AuthMiddleware(u *usecase.UserUsecase) fiber.Handler {
	var secretKey = utils.GetEnv("SECRET_JWT_KEY", "")
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		id, err := utils.DecodeToken(tokenString, secretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		_, err = u.GetById(id)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		return c.Next() // Continue if the token is valid
	}
}
