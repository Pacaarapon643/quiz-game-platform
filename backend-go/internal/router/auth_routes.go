package router

import (
	"quiz-game-backend/internal/cache"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/handlers"
	"quiz-game-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetAuthRoutes(api fiber.Router, db *database.Database, cfg *config.Config, redisClient *cache.RedisClient) {
	handler := handlers.NewAuthHandler(db, cfg, redisClient)

	auth := api.Group("/auth")

	// Public routes
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// Protected routes (ต้อง login ก่อน)
	auth.Post("/logout", middleware.Auth(cfg), handler.Logout)
	auth.Get("/me", middleware.Auth(cfg), handler.GetMe)
}
