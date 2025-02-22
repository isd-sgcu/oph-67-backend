package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/oph-67-backend/domain"
	"github.com/isd-sgcu/oph-67-backend/usecase"
	"github.com/isd-sgcu/oph-67-backend/utils"
)

func RoleMiddleware(u *usecase.UserUsecase, allowedRoles ...domain.Role) fiber.Handler {
	var secretKey = utils.GetEnv("SECRET_JWT_KEY", "")
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		id, err := utils.DecodeToken(tokenString, secretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		user, err := u.GetById(id)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}
		role := user.Role

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next() // Role is allowed, proceed to the next handler
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access forbidden: insufficient role permissions"})
	}
}
