package router

import (
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/handlers"
	"quiz-game-backend/internal/repository"
	"quiz-game-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetUserRoutes(api fiber.Router, db *database.Database) {
	// Create repository -> service -> handler chain
	repo := repository.NewUserResponse(db.DB)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	users := api.Group("/users")
	users.Get("/", handler.ListUsers)
	users.Get("/:id", handler.GetUserByID)
	users.Put("/:id", handler.UpdateProfile)
}
