package router

import (
	"quiz-game-backend/internal/cache"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, cfg *config.Config, db *database.Database, redisClient *cache.RedisClient) {
	// Health check
	app.Get("/health", handlers.HealthCheck(cfg, db))

	// API v1
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":     "Quiz Game API v1.0",
			"environment": cfg.Server.Environment,
		})
	})

	// Setup routes
	SetAuthRoutes(v1, db, cfg, redisClient) // Authentication routes
	SetUserRoutes(v1, db)                   // User routes

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
		})
	})
}
