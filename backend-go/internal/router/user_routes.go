package router

import (
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetUserRoutes(api fiber.Router, db *database.Database) {
	handler := handlers.NewUserHandler(db)

	users := api.Group("/users")
	users.Get("/", handler.ListUsers)
	users.Get("/:id", handler.GetUserByID)
	users.Put("/:id", handler.UpdateProfile)

}
