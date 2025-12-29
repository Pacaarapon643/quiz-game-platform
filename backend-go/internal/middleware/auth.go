package middleware

import (
	"log"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Auth middleware ตรวจสอบ JWT token จาก cookie
func Auth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง token จาก cookie
		token := utils.GetTokenFromCookie(c)

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - Missing token",
			})
		}

		// Validate token ด้วย RSA public key
		userID, err := utils.ExtractUserIDRSA(token, cfg.JWT.PublicKey)
		if err != nil {
			log.Println("Token validation error: ", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - Invalid token",
			})
		}

		// เก็บ user_id ใน context
		c.Locals("user_id", userID)
		return c.Next()
	}
}
